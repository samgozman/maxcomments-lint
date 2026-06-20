package maxcomments

import "testing"

func TestNolintForMaxcomments(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"bare nolint", "//nolint", true},
		{"nolint all", "//nolint:all", true},
		{"names maxcomments", "//nolint:maxcomments", true},
		{"maxcomments in a list", "//nolint:gocritic,maxcomments,godot", true},
		{"maxcomments with explanation", "//nolint:maxcomments // generated file", true},
		{"spaces around names", "//nolint: maxcomments , gocritic", true},
		{"other linters only", "//nolint:gocritic,godot", false},
		{"not a nolint", "// nolint is mentioned in prose", false},
		{"lookalike prefix", "//nolintfoo:maxcomments", false},
		{"plain comment", "// just a comment", false},
		{"block comment", "/* nolint:maxcomments */", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nolintForMaxcomments(tt.text); got != tt.want {
				t.Errorf("nolintForMaxcomments(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
