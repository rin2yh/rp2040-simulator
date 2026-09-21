// Package tabletest defines analyzers for table-driven test conventions.
package tabletest

import "golang.org/x/tools/go/analysis"

// Analyzer checks tabletest.Case types and repeated test cases.
var Analyzer = &analysis.Analyzer{
	Name: "tabletest",
	Doc:  "check table-driven test conventions",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	runNoFunc(pass)
	runRepeated(pass)
	return nil, nil
}
