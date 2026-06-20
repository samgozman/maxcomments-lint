package maxcomments

import (
	"go/ast"
	"go/token"
	"os"
	"strings"
)

// sourceLines holds a file's raw lines plus the set of line numbers that any
// comment occupies, so the ratio check can count "code lines": physical,
// non-blank lines that are not comment lines.
type sourceLines struct {
	lines       []string     // 1-based: line n is lines[n-1]
	commentLine map[int]bool // line numbers covered by a comment token
}

// newSourceLines reads the file backing the given AST from disk and records
// which lines are occupied by comments. Reading source is required because a
// token.FileSet knows line offsets but not whether a line is blank.
func newSourceLines(fset *token.FileSet, file *ast.File) (*sourceLines, error) {
	filename := fset.Position(file.Pos()).Filename

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	commentLine := make(map[int]bool)
	for _, group := range file.Comments {
		for _, c := range group.List {
			start := fset.Position(c.Pos()).Line
			end := fset.Position(c.End()).Line
			for ln := start; ln <= end; ln++ {
				commentLine[ln] = true
			}
		}
	}

	return &sourceLines{
		lines:       strings.Split(string(src), "\n"),
		commentLine: commentLine,
	}, nil
}

// codeLineCount returns how many lines in the inclusive [start, end] range are
// code lines: not blank and not occupied by a comment.
func (s *sourceLines) codeLineCount(start, end int) int {
	count := 0
	for ln := start; ln <= end; ln++ {
		if s.commentLine[ln] {
			continue
		}
		if ln < 1 || ln > len(s.lines) {
			continue
		}
		if strings.TrimSpace(s.lines[ln-1]) == "" {
			continue
		}
		count++
	}

	return count
}

// lineCount is the total number of lines in the source file.
func (s *sourceLines) lineCount() int {
	return len(s.lines)
}

// commentLineCount returns how many human-prose comment lines a group spans.
// Consecutive "//" lines each count as one line and multi-line "/* */" blocks
// count by their actual line span. Directive comments (see isDirective) are
// not counted.
func commentLineCount(fset *token.FileSet, group *ast.CommentGroup) int {
	if group == nil {
		return 0
	}

	total := 0
	for _, c := range group.List {
		if isDirective(c.Text) {
			continue
		}

		start := fset.Position(c.Pos()).Line
		end := fset.Position(c.End()).Line
		total += end - start + 1
	}

	return total
}

// isDirective reports whether a comment's text is a machine directive rather
// than human prose, and therefore should not count toward a comment budget.
//
// It recognizes Go's own directive convention (see go/ast: no space after the
// "//", a "tool:" prefix, or the "line " line directive) plus golangci-lint's
// "//nolint" suppression directives.
func isDirective(text string) bool {
	// Block comments (/* ... */) are never directives.
	rest, ok := strings.CutPrefix(text, "//")
	if !ok {
		return false
	}

	// golangci-lint suppression: "//nolint" or "//nolint:...".
	if rest == "nolint" || strings.HasPrefix(rest, "nolint:") {
		return true
	}

	// Go line directive: "//line file:line".
	if strings.HasPrefix(rest, "line ") {
		return true
	}

	// Go tool directive: "//word:..." where word is lowercase letters/digits
	// with no leading space, and something follows the colon.
	colon := strings.IndexByte(rest, ':')
	if colon <= 0 || colon+1 >= len(rest) {
		return false
	}
	for i := 0; i < colon; i++ {
		b := rest[i]
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9') {
			return false
		}
	}
	return true
}
