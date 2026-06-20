package maxcomments

import (
	"go/ast"
	"go/token"
	"strings"
)

// nolintForMaxcomments reports whether a comment is a golangci-lint nolint
// directive that suppresses this linter: a bare "//nolint", "//nolint:all",
// or a "//nolint:..." list that names "maxcomments" (or "all"). An optional
// "// explanation" trailing the linter list is ignored.
func nolintForMaxcomments(text string) bool {
	if !strings.HasPrefix(text, "//nolint") {
		return false
	}

	rest := strings.TrimPrefix(text, "//nolint")
	if rest == "" {
		return true // bare //nolint suppresses every linter
	}

	if !strings.HasPrefix(rest, ":") {
		return false // e.g. "//nolintfoo"
	}

	list := strings.TrimPrefix(rest, ":")
	if i := strings.Index(list, "//"); i >= 0 {
		list = list[:i] // drop the "// explanation" suffix
	}

	for _, name := range strings.Split(list, ",") {
		switch strings.TrimSpace(name) {
		case "all", "maxcomments":
			return true
		}
	}

	return false
}

// nolintInfo records where maxcomments nolint directives appear in a file.
type nolintInfo struct {
	// fileSuppressed is true when a nolint directive appears before the
	// package clause, suppressing the file-level checks.
	fileSuppressed bool
	// lines holds every line number carrying a maxcomments nolint directive,
	// used to suppress a function reported on that line.
	lines map[int]bool
}

// collectNolint scans the file's comments for maxcomments nolint directives.
func collectNolint(fset *token.FileSet, file *ast.File) nolintInfo {
	info := nolintInfo{lines: make(map[int]bool)}

	for _, group := range file.Comments {
		for _, c := range group.List {
			if !nolintForMaxcomments(c.Text) {
				continue
			}

			info.lines[fset.Position(c.Pos()).Line] = true
			if c.Pos() < file.Package {
				info.fileSuppressed = true
			}
		}
	}

	return info
}

// suppressesScope reports whether a function's diagnostics are suppressed by a
// nolint directive in its doc comment or on its signature line.
func (info nolintInfo) suppressesScope(fset *token.FileSet, scope *funcScope) bool {
	if groupHasNolint(scope.doc) {
		return true
	}

	return info.lines[fset.Position(scope.node.Pos()).Line]
}

func groupHasNolint(group *ast.CommentGroup) bool {
	if group == nil {
		return false
	}

	for _, c := range group.List {
		if nolintForMaxcomments(c.Text) {
			return true
		}
	}

	return false
}
