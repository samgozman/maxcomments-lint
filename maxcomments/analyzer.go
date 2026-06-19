// Package maxcomments implements a golangci-lint module plugin that caps
// the number of comment lines allowed inside a function or a file.
//
// The idea: heavily-commented code is sometimes a sign that the code itself
// needs simplifying rather than narrating. This linter does not judge
// comment quality -- only quantity -- so it pairs well with other linters
// that check comment style (godot, gocritic, revive).
package maxcomments

import (
	"go/ast"
	"go/token"

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
	total := 0

	if fn.Doc != nil {
		total += groupLineCount(pass.Fset, fn.Doc)
	}

	for _, group := range comments {
		if group == fn.Doc {
			continue
		}

		if group.Pos() >= fn.Pos() && group.End() <= fn.End() {
			total += groupLineCount(pass.Fset, group)
		}
	}

	if total > max {
		pass.Reportf(fn.Pos(), "function %q has %d comment lines, max allowed is %d",
			fn.Name.Name, total, max)
	}
}

// groupLineCount returns how many source lines a comment group spans,
// correctly handling both stacked "//" lines and multi-line "/* */" blocks.
func groupLineCount(fset *token.FileSet, group *ast.CommentGroup) int {
	if group == nil {
		return 0
	}

	start := fset.Position(group.Pos()).Line
	end := fset.Position(group.End()).Line

	return end - start + 1
}

func fileName(pass *analysis.Pass, file *ast.File) string {
	return pass.Fset.Position(file.Pos()).Filename
}
