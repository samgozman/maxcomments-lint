package maxcomments

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCodeLineCount_OutOfRange(t *testing.T) {
	s := &sourceLines{lines: []string{"package x"}, commentLine: map[int]bool{}}
	// The range runs past the only line; the out-of-range guard must skip the
	// missing lines rather than index out of bounds.
	if got := s.codeLineCount(1, 5); got != 1 {
		t.Fatalf("codeLineCount(1, 5) = %d, want 1", got)
	}
}

func TestNewSourceLines_ReadError(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "missing-7b21.go", "package x\n", parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// The file was parsed from memory and never written to disk, so reading it
	// back must fail.
	if _, err := newSourceLines(fset, file); err == nil {
		t.Fatal("expected read error for missing file, got nil")
	}
}

func TestRatioViolation(t *testing.T) {
	tests := []struct {
		name                            string
		comments, code, ratio, minLines int
		wantAllowed                     int
		wantViolated                    bool
	}{
		{"disabled when ratio zero", 5, 100, 0, 0, 0, false},
		{"within budget", 2, 30, 10, 0, 3, false},
		{"exactly at budget", 3, 30, 10, 0, 3, false},
		{"over budget", 4, 30, 10, 0, 3, true},
		{"floor division rounds down", 2, 29, 10, 0, 2, false},
		{"exempt below min lines", 5, 4, 3, 5, 0, false},
		{"checked at min lines", 5, 5, 3, 5, 1, true},
		{"zero code any comment violates", 1, 0, 10, 0, 0, true},
		{"zero code no comment ok", 0, 0, 10, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, violated := ratioViolation(tt.comments, tt.code, tt.ratio, tt.minLines)
			if allowed != tt.wantAllowed || violated != tt.wantViolated {
				t.Errorf("ratioViolation(%d,%d,%d,%d) = (%d,%v), want (%d,%v)",
					tt.comments, tt.code, tt.ratio, tt.minLines,
					allowed, violated, tt.wantAllowed, tt.wantViolated)
			}
		})
	}
}

func TestIsDirective(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"plain line comment", "// a normal sentence", false},
		{"plain word with colon and space after", "// note: this is prose", false},
		{"go generate", "//go:generate stringer -type=Foo", true},
		{"go embed", "//go:embed files", true},
		{"line directive", "//line foo.go:16", true},
		{"export directive", "//export MyFunc", false}, // cgo //export has a space, not a colon -> not a Go directive
		{"nolint bare", "//nolint", true},
		{"nolint with list", "//nolint:maxcomments,gocritic", true},
		{"nolint all", "//nolint:all", true},
		{"block comment never directive", "/* go:generate x */", false},
		{"colon but leading space invalidates", "// go:generate x", false},
		{"empty after slashes", "//", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDirective(tt.text); got != tt.want {
				t.Errorf("isDirective(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
