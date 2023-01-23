package golang

import (
	"go/ast"
	"io"
)

type Type struct {
	IsStruct bool
	Fields   []Field // If struct

	IsMap     bool
	KeyType   *Type // If map
	ValueType *Type // If map

	IsPointer bool
	IsArray   bool
	BaseType  *Type // If pointer or array

	IsIdent bool
	Ident   *Ident // If identifier

	Methods map[string]*ast.FuncDecl
}

func (t *Type) Packages() []string {
	if t == nil {
		return nil
	}
	pkgs := []string{}
	if t.IsIdent {
		pkgs = append(pkgs, t.Ident.From)
	}
	if t.IsMap {
		pkgs = append(pkgs, t.KeyType.Packages()...)
		pkgs = append(pkgs, t.ValueType.Packages()...)
	}
	if t.IsPointer || t.IsArray {
		pkgs = append(pkgs, t.BaseType.Packages()...)
	}
	return pkgs
}

func (t *Type) Write(io.Writer) error {
	panic("not implemented")
}
