package maxcomments_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/samgozman/maxcomments-lint/maxcomments"
)

func TestAnalyzer_FuncBudget(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "a")
}

func TestAnalyzer_NolintExclusion(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "b")
}

func TestAnalyzer_FuncLit(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "c")
}
