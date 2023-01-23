package golang

import (
	"go/ast"
)

func BuildType(currPkg string, importMap ImportMap, t ast.Expr) (*Type, error) {
	typ := &Type{
		Methods: map[string]*ast.FuncDecl{},
	}
	switch t := t.(type) {
	case *ast.Ident:
		typ.IsIdent = true
		typ.Ident = &Ident{
			Name: t.Name,
			From: importMap[t.Name],
		}
	case *ast.StarExpr:
		baseType, err := BuildType(currPkg, importMap, t.X)
		if err != nil {
			return nil, err
		}
		typ.IsPointer = true
		typ.BaseType = baseType
	case *ast.ArrayType:
		baseType, err := BuildType(currPkg, importMap, t.Elt)
		if err != nil {
			return nil, err
		}
		typ.IsArray = true
		typ.BaseType = baseType
	case *ast.MapType:
		keyType, err := BuildType(currPkg, importMap, t.Key)
		if err != nil {
			return nil, err
		}
		valueType, err := BuildType(currPkg, importMap, t.Value)
		if err != nil {
			return nil, err
		}
		typ.IsMap = true
		typ.KeyType = keyType
		typ.ValueType = valueType
	case *ast.StructType:
		fields := []Field{}
		for _, field := range t.Fields.List {
			fieldType, err := BuildType(currPkg, importMap, field.Type)
			if err != nil {
				return nil, err
			}
			fields = append(fields, Field{
				Name: field.Names[0].Name,
				Type: fieldType,
			})
		}
		typ.IsStruct = true
		typ.Fields = fields
	}
	return typ, nil
}
