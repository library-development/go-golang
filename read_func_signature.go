package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/library-development/go-nameconv"
)

// ReadFuncSignature reads the function signature of a function in a package.
// srcDir should be a directory with all your source code.
// Packages in srcDir should be in the format "github.com/username/repo".
func ReadFuncSignature(srcDir, pkg, funcName string) (*FuncSignature, error) {
	path := filepath.Join(srcDir, pkg)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	sig := FuncSignature{}
	for _, p := range pkgs {
		for _, f := range p.Files {
			imports, err := BuildImportMap(f)
			if err != nil {
				return nil, err
			}
			for _, d := range f.Decls {
				if fn, ok := d.(*ast.FuncDecl); ok {
					if fn.Recv == nil {
						if fn.Name.Name == funcName {
							for _, f := range fn.Type.Params.List {
								for _, n := range f.Names {
									typeID, err := BuildID(pkg, imports, f.Type)
									if err != nil {
										return nil, err
									}
									field := Field{
										Name: nameconv.Name{Words: []string{n.Name}},
										Type: typeID,
									}
									sig.Inputs = append(sig.Inputs, field)

								}
							}
							outID := 0
							for _, f := range fn.Type.Results.List {
								if len(f.Names) == 0 {
									outID++
									typeID, err := BuildID(pkg, imports, f.Type)
									if err != nil {
										return nil, err
									}
									field := Field{
										Name: nameconv.Name{Words: []string{"out" + strconv.Itoa(outID)}},
										Type: typeID,
									}
									sig.Outputs = append(sig.Outputs, field)
									continue
								}
								for range f.Names {
									outID++
									typeID, err := BuildID(pkg, imports, f.Type)
									if err != nil {
										return nil, err
									}
									field := Field{
										Name: nameconv.Name{Words: []string{"out" + strconv.Itoa(outID)}},
										Type: typeID,
									}
									sig.Outputs = append(sig.Outputs, field)
								}
							}
						}
					}
				}
			}
		}
	}
	return &sig, nil
}
