package maxcomments

import "testing"

func TestCompileIgnoreInvalid(t *testing.T) {
	if _, err := compileIgnore([]string{"valid", "("}); err == nil {
		t.Fatal("compileIgnore with an invalid regex should return an error")
	}
}

func TestMatchesAny(t *testing.T) {
	patterns, err := compileIgnore([]string{`vendor/`, `_test\.go$`})
	if err != nil {
		t.Fatalf("compileIgnore: %v", err)
	}

	tests := []struct {
		path string
		want bool
	}{
		{"/repo/vendor/foo/bar.go", true},
		{"/repo/pkg/thing_test.go", true},
		{"/repo/pkg/thing.go", false},
		{"/repo/testdata/x.go", false},
	}

	for _, tt := range tests {
		if got := matchesAny(patterns, tt.path); got != tt.want {
			t.Errorf("matchesAny(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestMatchesAnyEmpty(t *testing.T) {
	patterns, err := compileIgnore(nil)
	if err != nil {
		t.Fatalf("compileIgnore: %v", err)
	}
	if matchesAny(patterns, "/anything.go") {
		t.Error("empty pattern list should match nothing")
	}
}
