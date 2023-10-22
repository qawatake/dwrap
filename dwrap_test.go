package dwrap_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/qawatake/dwrap"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer_third_party(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, dwrap.NewAnalyzer("github.com/qawatake/derrors", "Wrap", "b/ignore"), "b/...")
}