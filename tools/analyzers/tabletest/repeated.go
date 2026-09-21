package tabletest

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const repeatedMessage = "repeated test cases can be expressed as a table-driven test"

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
		if !strings.HasSuffix(filepath.ToSlash(pass.Fset.Position(file.Pos()).Filename), "_test.go") {
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
	reported := make(map[token.Pos]bool)
	for start := 0; start < len(statements); start++ {
		for width := 1; width <= 3 && start+3*width <= len(statements); width++ {
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
			if count >= 3 && different && !reported[statements[start].Pos()] {
				pass.Reportf(statements[start].Pos(), repeatedMessage)
				reported[statements[start].Pos()] = true
			}
		}
	}
}

func describeGroup(pass *analysis.Pass, statements []ast.Stmt) statementGroup {
	group := statementGroup{
		normalized: normalizedStatements(pass, statements),
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

func normalizedStatements(pass *analysis.Pass, statements []ast.Stmt) string {
	type savedArguments struct {
		call *ast.CallExpr
		args []ast.Expr
	}
	var saved []savedArguments
	for _, statement := range statements {
		ast.Inspect(statement, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			obj, ok := calledObject(pass, call.Fun).(*types.Func)
			if !ok || isTestingAssertion(obj) {
				return true
			}
			saved = append(saved, savedArguments{call: call, args: call.Args})
			call.Args = []ast.Expr{ast.NewIdent("value")}
			return true
		})
	}
	formatted := formattedStatements(pass.Fset, statements)
	for _, item := range saved {
		item.call.Args = item.args
	}

	source := "package p\nfunc _() {\n" + formatted + "\n}"
	parsed, err := parser.ParseFile(token.NewFileSet(), "normalized.go", source, 0)
	if err != nil {
		return source
	}

	preserved := make(map[*ast.Ident]bool)
	ast.Inspect(parsed, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.SelectorExpr:
			preserved[node.Sel] = true
		case *ast.KeyValueExpr:
			if ident, ok := node.Key.(*ast.Ident); ok {
				preserved[ident] = true
			}
		}
		return true
	})
	ast.Inspect(parsed, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			if !preserved[node] {
				node.Name = "value"
			}
		case *ast.BasicLit:
			node.Value = normalizedLiteral(node.Kind)
		}
		return true
	})

	fn := parsed.Decls[0].(*ast.FuncDecl)
	return formattedStatements(token.NewFileSet(), fn.Body.List)
}

func normalizedLiteral(kind token.Token) string {
	switch kind {
	case token.INT:
		return "0"
	case token.FLOAT:
		return "0.0"
	case token.IMAG:
		return "0i"
	case token.CHAR:
		return "'x'"
	case token.STRING:
		return `""`
	default:
		return "0"
	}
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
