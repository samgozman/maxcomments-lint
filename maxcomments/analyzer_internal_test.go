package maxcomments

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// newPass builds a minimal analysis.Pass from in-memory source. The source is
// never written to disk, so a non-existent filename exercises the read-error
// paths.
func newPass(t *testing.T, filename, src string) (*analysis.Pass, *[]analysis.Diagnostic) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	var diags []analysis.Diagnostic
	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Report: func(d analysis.Diagnostic) { diags = append(diags, d) },
	}

	return pass, &diags
}

func TestRun_InvalidIgnorePattern(t *testing.T) {
	pass, _ := newPass(t, "x.go", "package x\n")
	if _, err := run(pass, Settings{Ignore: []string{"["}}); err == nil {
		t.Fatal("expected error for invalid ignore regex, got nil")
	}
}

func TestRun_SourceReadErrorPropagates(t *testing.T) {
	// A ratio setting forces a disk read of the non-existent file, so run must
	// propagate the error.
	pass, _ := newPass(t, "does-not-exist-9f3a.go", "package x\n\nfunc F() {}\n")
	if _, err := run(pass, Settings{File: FileSettings{Ratio: 1}}); err == nil {
		t.Fatal("expected error when source file cannot be read, got nil")
	}
}

func TestRun_FileHardCapReports(t *testing.T) {
	// Three comment lines over a budget of one; the hard cap needs no disk read.
	pass, diags := newPass(t, "x.go", "// one\n// two\n// three\npackage x\n")
	if _, err := run(pass, Settings{File: FileSettings{Lines: 1}}); err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(*diags) == 0 {
		t.Fatal("expected a file-level diagnostic, got none")
	}
}
