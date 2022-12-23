package golang

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/library-development/go-strutil"
)

// GenerateHandler generates ServeHTTP methods for each type declaration in the file.
func GenerateHTTPHandler(srcDir, pkg, typeName, outFile string) error {
	p, err := ParsePackage(srcDir, pkg)
	if err != nil {
		return err
	}
	pkgName := filepath.Base(pkg)
	pkgName = strings.TrimPrefix(pkgName, "go-")
	pkgName = strings.ToLower(pkgName)
	pkgName = strutil.Filter("abcdefghijklmnopqrstuvwxyz", pkgName)
	tmpl := `package {{ .PkgName }}

import (
	"encoding/json"
	"io"
	"net/http"
)
	
func (d *{{ .Name.Name }}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(d)
		return
	}
	if r.Method == http.MethodPost {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		var in struct {
			Method string
		}
		err = json.Unmarshal(b, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		switch in.Method {
		default:
			http.Error(w, "unknown method", http.StatusBadRequest)
		}
		return
	}
}
`
	fset := token.NewFileSet()
	dir, err := parser.ParseDir(fset, filepath.Join(srcDir, pkg), FilterTests, parser.ParseComments)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	offset := 0
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			p := n.End()
			buf.Write(b[offset:p])
			offset = int(p)
			t := template.Must(template.New("handler").Parse(tmpl))
			e := t.Execute(&buf, x)
			if e != nil {
				err = e
			}
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	buf.Write(b[offset:])
	err = os.WriteFile(file, buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	return RunGoimports(file)
}
