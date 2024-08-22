package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// test function
func TestNoOsExitInMain(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, NoOsExitInMainAnalyzer)
}
