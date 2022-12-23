package golang

import "go/ast"

type Package struct {
	Name  string
	Types map[string]*Type
	Funcs map[string]*ast.FuncDecl
}
