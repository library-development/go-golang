package golang

import "go/ast"

func BuildID(currentPkg string, importMap ImportMap, expr ast.Expr) (ID, error) {
	switch e := expr.(type) {
	case *ast.Ident:
		return ID{
			Pkg:  currentPkg,
			Name: e.Name,
		}, nil
	case *ast.SelectorExpr:
		pkg, _ := importMap.Resolve(e.X.(*ast.Ident).Name)
		return ID{
			Name: e.Sel.Name,
			Pkg:  pkg,
		}, nil
	}
	return ID{}, nil
}
