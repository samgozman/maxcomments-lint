package maxcomments

import (
	"go/ast"
	"go/token"
	"os"
	"strings"
)

// sourceLines holds a file's raw lines plus the line numbers any comment
// occupies, so the ratio check can count code lines (non-blank, non-comment).
type sourceLines struct {
	lines       []string     // 1-based: line n is lines[n-1]
	commentLine map[int]bool // line numbers covered by a comment token
}

// newSourceLines reads the file backing the AST from disk and records which
// lines comments occupy, so the ratio check can tell blank/comment from code.
func newSourceLines(fset *token.FileSet, file *ast.File) (*sourceLines, error) {
	filename := fset.Position(file.Pos()).Filename

	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(src), "\n")

	commentLine := make(map[int]bool)
	for _, group := range file.Comments {
		for _, c := range group.List {
			start := fset.Position(c.Pos())
			end := fset.Position(c.End()).Line

			// A trailing comment (preceded by code on its line) leaves that line
			// code; only its continuation lines, if any, are pure comment.
			first := start.Line
			if codeBefore(lines, start.Line, start.Column) {
				first++
			}
			for ln := first; ln <= end; ln++ {
				commentLine[ln] = true
			}
		}
	}

	return &sourceLines{
		lines:       lines,
		commentLine: commentLine,
	}, nil
}

// codeBefore reports whether non-whitespace precedes column col (1-based bytes)
// on the given line — i.e. a comment there is trailing, not whole-line.
func codeBefore(lines []string, line, col int) bool {
	if col < 2 || line < 1 || line > len(lines) {
		return false
	}

	prefix := lines[line-1]
	if col-1 < len(prefix) {
		prefix = prefix[:col-1]
	}

	return strings.TrimSpace(prefix) != ""
}

// codeLineCount returns how many lines in [start, end] are non-blank code.
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

func (s *sourceLines) lineCount() int {
	return len(s.lines)
}

// commentLineCount returns how many prose comment lines a group spans, by
// actual line span. Directive comments (see isDirective) are not counted.
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

// isDirective reports whether a comment is a machine directive rather than
// prose, so it is not counted. It recognizes Go's "//tool:..." and "//line "
// directives plus golangci-lint's "//nolint".
func isDirective(text string) bool {
	rest, ok := strings.CutPrefix(text, "//")
	if !ok {
		return false
	}

	if rest == "nolint" || strings.HasPrefix(rest, "nolint:") {
		return true
	}

	if strings.HasPrefix(rest, "line ") {
		return true
	}

	// Go tool directive "//word:x": a lowercase-alnum name, a colon, then a
	// lowercase-alnum char. That final check keeps "//note: foo" and "//https://"
	// (space or slash after the colon) as prose, matching go/ast.isDirective.
	colon := strings.IndexByte(rest, ':')
	if colon <= 0 || colon+1 >= len(rest) {
		return false
	}
	for i := 0; i <= colon+1; i++ {
		if i == colon {
			continue
		}
		if !isLowerAlnum(rest[i]) {
			return false
		}
	}
	return true
}

// isLowerAlnum reports whether b is a lowercase ASCII letter or a digit.
func isLowerAlnum(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= '0' && b <= '9'
}
