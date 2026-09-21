// Command vet runs repository-specific static analyzers.
package main

import (
	tabletestanalyzer "github.com/rin2yh/rp2040-simulator/tools/analyzers/tabletest"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		tabletestanalyzer.Analyzer,
	)
}
