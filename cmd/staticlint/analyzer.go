// Package provides an Analyzer which warns for direct calls to os.Exit in main function of main package.
package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// NoOsExitInMainAnalyzer is an analysis.Analyzer that checks and reports
// any direct calls to os.Exit in the main function of main package. It suggests
// to return an error instead of directly exiting the program.
var NoOsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noOsExitInMain",
	Doc:  "reports direct calls to os.Exit in main function of main package",
	Run:  run,
}

// run is the function implementing the NoOsExitInMainAnalyzer.
// It iterates over all the files in the package and reports direct calls to os.Exit.
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		// Check if it's a file in the main package
		if pass.Pkg.Name() != "main" {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			fun, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if fun.Sel.Name == "Exit" && qualifiedName(fun.X) == "os" {
				pass.Reportf(call.Pos(), "direct call to os.Exit found, consider returning an error instead")
			}

			return true
		})
	}

	return nil, nil
}

// qualifiedName checks if the Stmt has function Exit and returns its name.
// Returns an empty string if the Stmt is not a function or its name is not Exit.
func qualifiedName(x ast.Expr) string {
	ident, ok := x.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func main() {
	singlechecker.Main(NoOsExitInMainAnalyzer)
}
