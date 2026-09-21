package tabletestlinttest

import "github.com/rin2yh/rp2040-simulator/internal/tabletest"

type nested struct {
	Values map[string][]*callback
}

type callback func(error) bool

var valid = []tabletest.Case[string, struct{ Err bool }]{
	{Name: "valid", In: "input"},
}

var interfaceValue = tabletest.Case[string, interface{ Method() }]{}

var functionInput = []tabletest.Case[func(), bool]{ // want "tabletest.Case In type func\\(\\) contains a function type"
	{Name: "invalid"},
}

var nestedFunction = []tabletest.Case[string, nested]{ // want "tabletest.Case Want type github.com/rin2yh/rp2040-simulator/tabletestlinttest.nested contains a function type"
	{Name: "invalid"},
}
