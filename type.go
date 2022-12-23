package golang

import "go/ast"

type Type struct {
	IsStruct bool
	// If struct
	Fields []Field

	IsMap bool
	// If map
	KeyType   *Type
	ValueType *Type

	IsPointer bool
	IsArray   bool
	// If pointer or array
	BaseType *Type

	IsIdent bool
	// If identifier
	Ident *Ident

	Methods []*ast.FuncDecl
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
