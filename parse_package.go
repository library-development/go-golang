package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

func ParsePackage(srcDir, pkgName string) (*Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filepath.Join(srcDir, pkgName), FilterTests, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, ErrMultiplePackages
	}
	p := &Package{}

	// Find type declarations
	for _, pkg := range pkgs {
		p.Name = pkg.Name

		for _, f := range pkg.Files {
			importMap, err := BuildImportMap(f)
			if err != nil {
				return nil, err
			}
			for _, d := range f.Decls {
				switch d := d.(type) {
				case *ast.GenDecl:
					for _, s := range d.Specs {
						switch spec := s.(type) {
						case *ast.TypeSpec:
							typ, err := BuildType(pkgName, importMap, spec.Type)
							if err != nil {
								return nil, err
							}
							p.Types[spec.Name.Name] = typ
						}
					}
				}
			}
		}
	}

	// Find func/method declarations
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			importMap, err := BuildImportMap(f)
			if err != nil {
				return nil, err
			}
			for _, d := range f.Decls {
				switch d := d.(type) {
				case *ast.FuncDecl:
					if d.Recv != nil {
						// Method
						if len(d.Recv.List) != 1 {
							return nil, ErrMultipleReceivers
						}
						ident, err := BuildIdent(pkgName, importMap, d.Recv.List[0].Type)
						if err != nil {
							return nil, err
						}
						p.Types[ident.Name].Methods = append(p.Types[ident.Name].Methods, d)
					} else {
						p.Funcs[d.Name.Name] = d
					}
				}
			}
		}
	}

	return p, nil
}
