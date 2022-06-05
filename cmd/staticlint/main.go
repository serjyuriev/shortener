package main

import (
	"go/ast"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

// ExitCheckAnalyzer checks package main function
// for direct calls to os.Exit().
var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for direct calls to os.Exit() in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			ast.Inspect(file, func(node ast.Node) bool {
				if x, ok := node.(*ast.CallExpr); ok {
					if f, ok := x.Fun.(*ast.SelectorExpr); ok {
						if p, ok := f.X.(*ast.Ident); ok {
							if p.Name == "os" && f.Sel.Name == "Exit" {
								pass.Reportf(f.Pos(), "os.Exit() call in main")
							}
						}
					}
				}
				return true
			})
		}
	}
	return nil, nil
}

func main() {
	mychecks := []*analysis.Analyzer{
		ExitCheckAnalyzer,       // os.Exit() calls in main function
		errcheck.Analyzer,       // unchecked errors
		ineffassign.Analyzer,    // inefficient assign
		assign.Analyzer,         // useless assignments
		bools.Analyzer,          // mistakes with boolean operators
		copylock.Analyzer,       // locks that are being passed by value
		errorsas.Analyzer,       // checks whether in errors.As second argument is a pointer to type that implements error interface
		fieldalignment.Analyzer, // rearrange struct fields for less memory consumption
		httpresponse.Analyzer,   // use predefinied http response codes
		lostcancel.Analyzer,     // always call a context cancellation function
		nilfunc.Analyzer,        // do not compare with nil
		printf.Analyzer,         // right formats in printf
		shadow.Analyzer,         // search for shadowed variables
		unmarshal.Analyzer,      // passing non-pointer or non-interface types to unmarshal and decode functions
		unreachable.Analyzer,    // search for unreachable code
		unusedresult.Analyzer,   // search for unused result
		unusedwrite.Analyzer,    // search for unused write to variables
	}

	checks := map[string]bool{
		"SA":     true,
		"S1028":  true,
		"ST1006": true,
		"QF1001": true,
	}
	for _, v := range staticcheck.Analyzers {
		if checks[v.Name] {
			mychecks = append(mychecks, v)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
