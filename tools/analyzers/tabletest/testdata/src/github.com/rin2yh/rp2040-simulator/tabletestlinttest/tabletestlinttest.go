package tabletestlinttest

import "github.com/rin2yh/rp2040-simulator/internal/tabletest"

type scalarInput string

type scalarWant struct {
	Err  bool
	Code int
}

type callback func(error) bool

type functionField struct {
	Check callback
}

type pointerToFunction struct {
	Check *callback
}

type sliceOfFunctions struct {
	Checks []callback
}

type mapOfFunctions struct {
	Checks map[string]callback
}

var validScalar = tabletest.Case[scalarInput, bool]{}
var validStruct = tabletest.Case[struct{ Value int }, scalarWant]{}

var functionInput = tabletest.Case[func(), bool]{}                                 // want "tabletest.Case In type func\\(\\) contains a function type"
var functionWant = tabletest.Case[string, func(error) bool]{}                      // want "tabletest.Case Want type func\\(error\\) bool contains a function type"
var functionInStruct = tabletest.Case[string, functionField]{}                     // want "tabletest.Case Want type github.com/rin2yh/rp2040-simulator/tabletestlinttest.functionField contains a function type"
var functionThroughPointer = tabletest.Case[pointerToFunction, bool]{}             // want "tabletest.Case In type github.com/rin2yh/rp2040-simulator/tabletestlinttest.pointerToFunction contains a function type"
var functionThroughSlice = tabletest.Case[string, sliceOfFunctions]{}              // want "tabletest.Case Want type github.com/rin2yh/rp2040-simulator/tabletestlinttest.sliceOfFunctions contains a function type"
var functionThroughMap = tabletest.Case[string, mapOfFunctions]{}                  // want "tabletest.Case Want type github.com/rin2yh/rp2040-simulator/tabletestlinttest.mapOfFunctions contains a function type"
var functionInInterface = tabletest.Case[string, interface{ Check(error) bool }]{} // want "tabletest.Case Want type interface{Check\\(error\\) bool} contains a function type"
