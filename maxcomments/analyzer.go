// Package maxcomments implements a golangci-lint module plugin that caps
// the number of comment lines allowed inside a function or a file.
//
// The idea: heavily-commented code is sometimes a sign that the code itself
// needs simplifying rather than narrating. This linter does not judge
// comment quality -- only quantity -- so it pairs well with other linters
// that check comment style (godot, gocritic, revive).
package maxcomments //nolint:maxcomments

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Settings configures the maxcomments analyzer. Every field is optional; a
// value of 0 (or, for Ignore, an empty list) disables that particular check.
// Per-scope budgets are grouped under Func and File.
type Settings struct {
	// Func holds the per-function comment budgets.
	Func FuncSettings `json:"func"`

	// File holds the per-file comment budgets.
	File FileSettings `json:"file"`

	// RatioMinLines suppresses the ratio checks for any scope with fewer than
	// this many code lines, so small functions are not penalised. 0 means no
	// floor (the ratio applies to every scope). It is scope-independent and so
	// lives outside Func/File.
	RatioMinLines int `json:"ratio-min-lines"`

	// Ignore is a list of regular expressions matched against each file's
	// path. A file whose path matches any pattern is skipped entirely. An
	// empty list checks every file.
	Ignore []string `json:"ignore"`

	// CheckGenerated controls whether machine-generated files are checked.
	// Generated files (those carrying the standard
	// "// Code generated ... DO NOT EDIT." marker) are skipped by default;
	// set this to true to check them like any other file.
	CheckGenerated bool `json:"check-generated"`
}

// FuncSettings holds the per-function comment budgets. Doc comments (the block
// directly above `func`) and body comments (everything inside the body) are
// tracked separately so a long doc comment cannot trip the body budget, and
// vice versa.
type FuncSettings struct {
	// BodyLines is the maximum number of body comment lines allowed inside a
	// single function, excluding its doc comment. 0 disables the check.
	BodyLines int `json:"body-lines"`

	// DocLines is the maximum number of doc comment lines allowed on a single
	// function: the comment block directly above its `func` keyword. 0 disables
	// the check.
	DocLines int `json:"doc-lines"`

	// Ratio enables the per-function body-comments-to-code ratio check: at most
	// one body comment line is allowed per Ratio code lines (so the allowed
	// budget is floor(codeLines / Ratio)). Doc comments are not included.
	// 0 disables it.
	Ratio int `json:"ratio"`
}

// FileSettings holds the per-file comment budgets, which count every comment
// group in the file (including the package doc comment).
type FileSettings struct {
	// Lines is the maximum number of comment lines allowed in a single file.
	// 0 disables the check.
	Lines int `json:"lines"`

	// Ratio enables the file-scope comments-to-code ratio check: at most one
	// comment line per Ratio code lines in the file. 0 disables it.
	Ratio int `json:"ratio"`
}

// NewAnalyzer builds the comment-budget analyzer for the given settings.
func NewAnalyzer(settings Settings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "maxcomments",
		Doc:  "reports functions or files whose comments exceed a configured line budget",
		URL:  "https://github.com/samgozman/maxcomments-lint",
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, settings)
		},
	}
}

func run(pass *analysis.Pass, settings Settings) (any, error) {
	ignore, err := compileIgnore(settings.Ignore)
	if err != nil {
		return nil, err
	}

	for _, file := range pass.Files {
		if matchesAny(ignore, fileName(pass, file)) {
			continue
		}

		if !settings.CheckGenerated && ast.IsGenerated(file) {
			continue
		}

		if err := checkFile(pass, file, settings); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File, settings Settings) error {
	// Reading source is only needed by the ratio checks, so do it lazily.
	var src *sourceLines
	if settings.Func.Ratio > 0 || settings.File.Ratio > 0 {
		s, err := newSourceLines(pass.Fset, file)
		if err != nil {
			return err
		}
		src = s
	}

	nolint := collectNolint(pass.Fset, file)

	if !nolint.fileSuppressed {
		checkFileBudget(pass, file, settings, src)
	}

	if settings.Func.BodyLines <= 0 && settings.Func.Ratio <= 0 && settings.Func.DocLines <= 0 {
		return nil
	}

	for _, scope := range collectFuncScopes(file) {
		if nolint.suppressesScope(pass.Fset, scope) {
			continue
		}
		checkScope(pass, scope, settings, src)
	}

	return nil
}

func checkFileBudget(pass *analysis.Pass, file *ast.File, settings Settings, src *sourceLines) {
	total := 0
	for _, group := range file.Comments {
		total += commentLineCount(pass.Fset, group)
	}

	if settings.File.Lines > 0 && total > settings.File.Lines {
		pass.Reportf(file.Package, "file %q has %d comment lines, max allowed is %d (file.lines)",
			fileName(pass, file), total, settings.File.Lines)
	}

	if settings.File.Ratio > 0 && src != nil {
		code := src.codeLineCount(1, src.lineCount())
		if allowed, violated := ratioViolation(total, code, settings.File.Ratio, settings.RatioMinLines); violated {
			pass.Reportf(file.Package,
				"file %q has %d comment lines for %d code lines, max allowed is %d (file.ratio)",
				fileName(pass, file), total, code, allowed)
		}
	}
}

func checkScope(pass *analysis.Pass, scope *funcScope, settings Settings, src *sourceLines) {
	doc := commentLineCount(pass.Fset, scope.doc)
	body := 0
	for _, group := range scope.comments {
		body += commentLineCount(pass.Fset, group)
	}

	if settings.Func.DocLines > 0 && doc > settings.Func.DocLines {
		pass.Reportf(scope.node.Pos(), "%s has %d doc comment lines, max allowed is %d (func.doc-lines)",
			scope.name, doc, settings.Func.DocLines)
	}

	if settings.Func.BodyLines > 0 && body > settings.Func.BodyLines {
		pass.Reportf(scope.node.Pos(), "%s has %d body comment lines, max allowed is %d (func.body-lines)",
			scope.name, body, settings.Func.BodyLines)
	}

	if settings.Func.Ratio > 0 && src != nil {
		start := pass.Fset.Position(scope.node.Pos()).Line
		end := pass.Fset.Position(scope.node.End()).Line
		code := src.codeLineCount(start, end)
		if allowed, violated := ratioViolation(body, code, settings.Func.Ratio, settings.RatioMinLines); violated {
			pass.Reportf(scope.node.Pos(),
				"%s has %d body comment lines for %d code lines, max allowed is %d (func.ratio)",
				scope.name, body, code, allowed)
		}
	}
}

// ratioViolation reports whether commentLines exceeds the budget implied by
// allowing one comment line per `ratio` code lines. Scopes with fewer than
// minLines code lines are exempt. The returned allowed value is the budget.
func ratioViolation(commentLines, codeLines, ratio, minLines int) (allowed int, violated bool) {
	if ratio <= 0 || codeLines < minLines {
		return 0, false
	}

	allowed = codeLines / ratio
	return allowed, commentLines > allowed
}

func fileName(pass *analysis.Pass, file *ast.File) string {
	return pass.Fset.Position(file.Pos()).Filename
}
