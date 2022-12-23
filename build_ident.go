package golang

import "go/ast"

func BuildIdent(currentPkg string, importMap ImportMap, expr ast.Expr) (*Ident, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		return &Ident{
			From: currentPkg,
			Name: e.Name,
		}, nil
	case *ast.SelectorExpr:
		pkg, _ := importMap.Resolve(e.X.(*ast.Ident).Name)
		return &Ident{
			From: pkg,
			Name: e.Sel.Name,
		}, nil
	case *ast.StarExpr:
		return BuildIdent(currentPkg, importMap, e.X)
	case *ast.ArrayType:
		return BuildIdent(currentPkg, importMap, e.Elt)
	case *ast.MapType:
		return BuildIdent(currentPkg, importMap, e.Value)
	}
	return nil, nil
}
