package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var NoOsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noOsExitInMain",
	Doc:  "reports direct calls to os.Exit in main function of main package",
	Run:  run,
}

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
