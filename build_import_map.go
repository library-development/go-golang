package golang

import "go/ast"

func BuildImportMap(f *ast.File) (ImportMap, error) {
	m := ImportMap{}
	for _, s := range f.Imports {
		if s.Name != nil {
			m[s.Path.Value] = s.Name.Name
		}
		m.AddPackage(s.Path.Value)
	}
	return m, nil
}
