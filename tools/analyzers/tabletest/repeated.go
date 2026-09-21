package tabletest

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const repeatedMessage = "repeated test cases can be expressed as a table-driven test"
const maxStatementGroup = 3

type statementGroup struct {
	normalized string
	original   string
	calls      []types.Object
	assertion  bool
	subject    bool
	ordered    bool
}

func runRepeated(pass *analysis.Pass) {
	for _, file := range pass.Files {
		if !strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, "_test.go") {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || fn.Body == nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				if block, ok := node.(*ast.BlockStmt); ok {
					reportRepeatedGroups(pass, block.List)
				}
				return true
			})
		}
	}
}

func reportRepeatedGroups(pass *analysis.Pass, statements []ast.Stmt) {
	for start := 0; start < len(statements); start++ {
		for width := 1; width <= maxStatementGroup && start+3*width <= len(statements); width++ {
			first := describeGroup(pass, statements[start:start+width])
			if !first.assertion || !first.subject || first.ordered {
				continue
			}

			count := 1
			different := false
			for next := start + width; next+width <= len(statements); next += width {
				candidate := describeGroup(pass, statements[next:next+width])
				if candidate.normalized != first.normalized || !slices.Equal(candidate.calls, first.calls) {
					break
				}
				count++
				different = different || candidate.original != first.original
			}
			if count >= 3 && different {
				pass.Reportf(statements[start].Pos(), repeatedMessage)
				break
			}
		}
	}
}

func describeGroup(pass *analysis.Pass, statements []ast.Stmt) statementGroup {
	group := statementGroup{
		normalized: structuralStatements(pass.Fset, statements),
		original:   formattedStatements(pass.Fset, statements),
		ordered:    hasBareSubjectCall(pass, statements),
	}
	for _, statement := range statements {
		ast.Inspect(statement, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			obj := calledObject(pass, call.Fun)
			if obj == nil {
				return true
			}
			fn, ok := obj.(*types.Func)
			if !ok {
				return true
			}
			group.calls = append(group.calls, fn)
			if isTestingAssertion(fn) {
				group.assertion = true
			} else {
				group.subject = true
			}
			return true
		})
	}
	return group
}

func hasBareSubjectCall(pass *analysis.Pass, statements []ast.Stmt) bool {
	for _, statement := range statements {
		expression, ok := statement.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := expression.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		if obj := calledObject(pass, call.Fun); obj != nil && !isTestingAssertion(obj) {
			if _, ok := obj.(*types.Func); ok {
				return true
			}
		}
	}
	return false
}

func calledObject(pass *analysis.Pass, expression ast.Expr) types.Object {
	switch expression := expression.(type) {
	case *ast.Ident:
		return pass.TypesInfo.Uses[expression]
	case *ast.SelectorExpr:
		if selection := pass.TypesInfo.Selections[expression]; selection != nil {
			return selection.Obj()
		}
		return pass.TypesInfo.Uses[expression.Sel]
	case *ast.IndexExpr:
		return calledObject(pass, expression.X)
	case *ast.IndexListExpr:
		return calledObject(pass, expression.X)
	case *ast.ParenExpr:
		return calledObject(pass, expression.X)
	default:
		return nil
	}
}

func isTestingAssertion(obj types.Object) bool {
	fn, ok := obj.(*types.Func)
	if !ok || fn.Pkg() == nil || fn.Pkg().Path() != "testing" {
		return false
	}
	switch fn.Name() {
	case "Error", "Errorf", "Fatal", "Fatalf":
		return true
	default:
		return false
	}
}

func structuralStatements(fset *token.FileSet, statements []ast.Stmt) string {
	filter := func(name string, value reflect.Value) bool {
		if name == "Name" || name == "Value" || name == "Obj" || name == "Args" {
			return false
		}
		return value.Type() != reflect.TypeFor[token.Pos]()
	}

	var buf bytes.Buffer
	for _, statement := range statements {
		ast.Fprint(&buf, fset, statement, filter)
	}
	return buf.String()
}

func formattedStatements(fset *token.FileSet, statements []ast.Stmt) string {
	var buf bytes.Buffer
	for _, statement := range statements {
		if err := format.Node(&buf, fset, statement); err != nil {
			return ""
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
