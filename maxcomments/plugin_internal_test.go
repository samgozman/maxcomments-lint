package maxcomments

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestNew(t *testing.T) {
	p, err := New(map[string]any{"max-func-lines": 5})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers: %v", err)
	}
	if len(analyzers) != 1 {
		t.Fatalf("BuildAnalyzers returned %d analyzers, want 1", len(analyzers))
	}
	if got := analyzers[0].Name; got != "maxcomments" {
		t.Errorf("analyzer name = %q, want %q", got, "maxcomments")
	}

	if got := p.GetLoadMode(); got != register.LoadModeSyntax {
		t.Errorf("GetLoadMode = %q, want %q", got, register.LoadModeSyntax)
	}
}

func TestNew_InvalidSettings(t *testing.T) {
	// DecodeSettings JSON-decodes into Settings; a scalar cannot decode into
	// the struct, so New must surface the error rather than panic.
	if _, err := New(12345); err == nil {
		t.Fatal("expected error decoding invalid settings, got nil")
	}
}
