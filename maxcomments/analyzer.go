// Package maxcomments implements a golangci-lint module plugin that caps
// the number of comment lines allowed inside a function or a file.
//
// The idea: heavily-commented code is sometimes a sign that the code itself
// needs simplifying rather than narrating. This linter does not judge
// comment quality -- only quantity -- so it pairs well with other linters
// that check comment style (godot, gocritic, revive).
package maxcomments

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Settings configures the maxcomments analyzer. Both fields are optional;
// a value of 0 disables that particular check.
type Settings struct {
	// MaxFuncLines is the maximum number of comment lines allowed inside a
	// single function, counting its doc comment plus any comments in its
	// body. 0 disables the check.
	MaxFuncLines int `json:"max-func-lines"`

	// MaxFileLines is the maximum number of comment lines allowed in a
	// single file, counting every comment group in the file. 0 disables
	// the check.
	MaxFileLines int `json:"max-file-lines"`
}

// NewAnalyzer builds the comment-budget analyzer for the given settings.
func NewAnalyzer(settings Settings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "maxcomments",
		Doc:  "reports functions or files whose comments exceed a configured line budget",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, settings)
		},
	}
}

func run(pass *analysis.Pass, settings Settings) (interface{}, error) {
	for _, file := range pass.Files {
		checkFile(pass, file, settings)
	}

	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File, settings Settings) {
	if settings.MaxFileLines > 0 {
		total := 0
		for _, group := range file.Comments {
			total += groupLineCount(pass.Fset, group)
		}

		if total > settings.MaxFileLines {
			pass.Reportf(file.Package, "file %q has %d comment lines, max allowed is %d",
				fileName(pass, file), total, settings.MaxFileLines)
		}
	}

	if settings.MaxFuncLines <= 0 {
		return
	}

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		checkFunc(pass, fn, file.Comments, settings.MaxFuncLines)
	}
}

func checkFunc(pass *analysis.Pass, fn *ast.FuncDecl, comments []*ast.CommentGroup, max int) {
	checkFuncNode(pass, fn.Pos(), fn.End(), fn.Doc,
		fmt.Sprintf("function %q", fn.Name.Name), fn.Body, comments, max)
}

// checkFuncNode checks the comment budget for a single function scope — either
// a top-level FuncDecl or an anonymous FuncLit.  It excludes comments that
// belong to direct child FuncLit nodes (those are checked recursively).
func checkFuncNode(
	pass *analysis.Pass,
	nodePos, nodeEnd token.Pos,
	docGroup *ast.CommentGroup,
	name string,
	body *ast.BlockStmt,
	comments []*ast.CommentGroup,
	max int,
) {
	// Collect FuncLit children that are direct descendants of this scope,
	// i.e. not nested inside another FuncLit.  We stop descending into any
	// FuncLit we find so that deeper ones are handled by the recursive call.
	var directLits []*ast.FuncLit
	if body != nil {
		ast.Inspect(body, func(n ast.Node) bool {
			if n == body {
				return true
			}
			if lit, ok := n.(*ast.FuncLit); ok {
				directLits = append(directLits, lit)
				return false // don't recurse into the literal
			}
			return true
		})
	}

	total := 0
	if docGroup != nil {
		total += groupLineCount(pass.Fset, docGroup)
	}

	for _, group := range comments {
		if group == docGroup {
			continue
		}
		if group.Pos() < nodePos || group.End() > nodeEnd {
			continue
		}
		if insideAnyFuncLit(group.Pos(), group.End(), directLits) {
			continue
		}
		total += groupLineCount(pass.Fset, group)
	}

	if total > max {
		pass.Reportf(nodePos, "%s has %d comment lines, max allowed is %d", name, total, max)
	}

	// Recursively check each direct child FuncLit.
	for _, lit := range directLits {
		checkFuncNode(pass, lit.Pos(), lit.End(), nil,
			"function literal", lit.Body, comments, max)
	}
}

// insideAnyFuncLit reports whether the range [pos, end] falls inside any of
// the given function literals.
func insideAnyFuncLit(pos, end token.Pos, lits []*ast.FuncLit) bool {
	for _, lit := range lits {
		if pos >= lit.Pos() && end <= lit.End() {
			return true
		}
	}
	return false
}

// groupLineCount returns the number of non-nolint comment lines in group,
// correctly handling both stacked "//" lines and multi-line "/* */" blocks.
func groupLineCount(fset *token.FileSet, group *ast.CommentGroup) int {
	if group == nil {
		return 0
	}

	count := 0
	for _, c := range group.List {
		if isNolintComment(c) {
			continue
		}
		start := fset.Position(c.Pos()).Line
		end := fset.Position(c.End()).Line
		count += end - start + 1
	}
	return count
}

// isNolintComment reports whether c is a //nolint directive.  Such lines are
// excluded from the comment-line budget because they are tooling directives,
// not documentation or narrative comments.
//
// Only single-line // comments are checked: golangci-lint's nolint mechanism
// requires the //nolint form, so /* nolint */ block comments are not a concern.
func isNolintComment(c *ast.Comment) bool {
	text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
	return strings.HasPrefix(text, "nolint")
}

func fileName(pass *analysis.Pass, file *ast.File) string {
	return pass.Fset.Position(file.Pos()).Filename
}
