package maxcomments

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("maxcomments", New)
}

// New builds the plugin instance for golangci-lint's module plugin system.
// It is wired up via register.Plugin in init() above.
func New(conf any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](conf)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: settings}, nil
}

type plugin struct {
	settings Settings
}

var _ register.LinterPlugin = &plugin{}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{NewAnalyzer(p.settings)}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
