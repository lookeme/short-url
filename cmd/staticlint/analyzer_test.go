package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// Write your test function
func TestNoOsExitInMain(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, NoOsExitInMainAnalyzer)
}
