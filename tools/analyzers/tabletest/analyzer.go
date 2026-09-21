// Package tabletest defines an analyzer for tabletest.Case type arguments.
package tabletest

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const tabletestPackage = "github.com/rin2yh/rp2040-simulator/internal/tabletest"

// Analyzer rejects function types in tabletest.Case input and expected values.
var Analyzer = &analysis.Analyzer{
	Name: "tabletest",
	Doc:  "reject function types in tabletest.Case input and expected values",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for ident, instance := range pass.TypesInfo.Instances {
		obj, ok := pass.TypesInfo.Uses[ident].(*types.TypeName)
		if !ok || obj.Pkg() == nil || obj.Pkg().Path() != tabletestPackage || obj.Name() != "Case" {
			continue
		}

		for i, name := range []string{"In", "Want"} {
			typeArg := instance.TypeArgs.At(i)
			if containsFunction(typeArg, make(map[types.Type]bool)) {
				pass.Reportf(ident.Pos(), "tabletest.Case %s type %s contains a function type", name, typeArg)
			}
		}
	}
	return nil, nil
}

func containsFunction(t types.Type, seen map[types.Type]bool) bool {
	t = types.Unalias(t)
	if seen[t] {
		return false
	}
	seen[t] = true

	switch t := t.(type) {
	case *types.Signature:
		return true
	case *types.Named:
		return containsFunction(t.Underlying(), seen)
	case *types.Pointer:
		return containsFunction(t.Elem(), seen)
	case *types.Slice:
		return containsFunction(t.Elem(), seen)
	case *types.Array:
		return containsFunction(t.Elem(), seen)
	case *types.Map:
		return containsFunction(t.Key(), seen) || containsFunction(t.Elem(), seen)
	case *types.Chan:
		return containsFunction(t.Elem(), seen)
	case *types.Struct:
		for i := range t.NumFields() {
			if containsFunction(t.Field(i).Type(), seen) {
				return true
			}
		}
	case *types.Interface:
		for i := range t.NumExplicitMethods() {
			if containsFunction(t.ExplicitMethod(i).Type(), seen) {
				return true
			}
		}
		for i := range t.NumEmbeddeds() {
			if containsFunction(t.EmbeddedType(i), seen) {
				return true
			}
		}
	case *types.Tuple:
		for i := range t.Len() {
			if containsFunction(t.At(i).Type(), seen) {
				return true
			}
		}
	case *types.TypeParam:
		return containsFunction(t.Constraint(), seen)
	case *types.Union:
		for i := range t.Len() {
			if containsFunction(t.Term(i).Type(), seen) {
				return true
			}
		}
	}
	return false
}
