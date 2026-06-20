package maxcomments_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/samgozman/maxcomments-lint/maxcomments"
)

func TestAnalyzer_FuncBudget(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "funclines")
}

func TestAnalyzer_ExcludesDirectives(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "directives")
}

func TestAnalyzer_ChecksClosuresIndependently(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "closures")
}

func TestAnalyzer_FuncRatio(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncRatio: 3})
	analysistest.Run(t, testdata, analyzer, "funcratio")
}

func TestAnalyzer_FileRatio(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFileRatio: 3})
	analysistest.Run(t, testdata, analyzer, "fileratio")
}

func TestAnalyzer_RatioMinLines(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncRatio: 3, RatioMinLines: 5})
	analysistest.Run(t, testdata, analyzer, "ratiomin")
}

func TestAnalyzer_NolintFile(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFileLines: 2, MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "nolintfile")
}

func TestAnalyzer_NolintFunc(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{MaxFuncLines: 2})
	analysistest.Run(t, testdata, analyzer, "nolintfunc")
}

func TestAnalyzer_IgnorePaths(t *testing.T) {
	testdata := analysistest.TestData()
	analyzer := maxcomments.NewAnalyzer(maxcomments.Settings{
		MaxFuncLines: 2,
		Ignore:       []string{`ignored\.go$`},
	})
	analysistest.Run(t, testdata, analyzer, "ignore")
}
