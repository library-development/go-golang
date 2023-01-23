package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type Workdir string

// CloneGithubRepo clones the given repo into the workdir at workdir/org/repo.
func (dir Workdir) CloneGithubRepo(org, repo string) error {
	gitURL := fmt.Sprintf("https://github.com/%s/%s.git", org, repo)

	// Make sure the org directory exists.
	err := os.MkdirAll(filepath.Join(string(dir), org), os.ModePerm)
	if err != nil {
		return fmt.Errorf("MkdirAll: %s", err)
	}

	// Clone the repo.
	cmd := exec.Command("git", "clone", gitURL)
	cmd.Dir = filepath.Join(string(dir), org)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone: %s: %s", err, out)
	}

	// Run go work use.
	cmd = exec.Command("go", "work", "use", filepath.Join(org, repo))
	cmd.Dir = string(dir)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go work use: %s: %s", err, out)
	}

	return nil
}

// Build runs go build -o output target in the workdir.
func (dir Workdir) Build(path, output string) error {
	cmd := exec.Command("go", "build", "-o", output, path)
	cmd.Dir = string(dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build: %s: %s", err, out)
	}
	return nil
}

func (dir Workdir) pkgPath(p string) (string, error) {
	goworkfile := filepath.Join(string(dir), "go.work")
	b, err := os.ReadFile(goworkfile)
	if err != nil {
		return "", fmt.Errorf("ReadFile: %s", err)
	}
	workfile, err := modfile.ParseWork(goworkfile, b, nil)
	if err != nil {
		return "", fmt.Errorf("ParseWork: %s", err)
	}
	for _, pkg := range workfile.Use {
		gomodfile := filepath.Join(string(dir), pkg.Path, "go.mod")
		b, err := os.ReadFile(gomodfile)
		if err != nil {
			return "", fmt.Errorf("ReadFile: %s", err)
		}
		mf, err := modfile.Parse(gomodfile, b, nil)
		if err != nil {
			return "", fmt.Errorf("go.mod parse: %s", err)
		}
		if mf.Module.Mod.Path != p {
			continue
		}
		return filepath.Join(string(dir), pkg.Path), nil
	}
	return "", fmt.Errorf("package %s not found", p)
}

func (dir Workdir) ParsePackage(path string) (*Package, error) {
	pkgPath, err := dir.pkgPath(path)
	if err != nil {
		return nil, fmt.Errorf("pkgPath: %s", err)
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, FilterTests, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, ErrMultiplePackages
	}
	p := &Package{
		Types: map[string]*Type{},
		Funcs: map[string]*ast.FuncDecl{},
	}
	for _, pkg := range pkgs {
		p.Name = pkg.Name
		// Find type declarations
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
							typ, err := BuildType(pkg.Name, importMap, spec.Type)
							if err != nil {
								return nil, err
							}
							p.Types[spec.Name.Name] = typ
						}
					}
				}
			}
		}
		// Find func/method declarations
		for _, f := range pkg.Files {
			importMap, err := BuildImportMap(f)
			if err != nil {
				return nil, err
			}
			for _, d := range f.Decls {
				switch d := d.(type) {
				case *ast.FuncDecl:
					if d.Recv != nil { // Method
						if len(d.Recv.List) != 1 {
							return nil, ErrMultipleReceivers
						}
						ident, err := BuildIdent(pkg.Name, importMap, d.Recv.List[0].Type)
						if err != nil {
							return nil, err
						}
						t, ok := p.Types[ident.Name]
						if !ok {
							return nil, fmt.Errorf("type %s not found", ident.Name)
						}
						t.Methods[d.Name.Name] = d
					} else {
						p.Funcs[d.Name.Name] = d
					}
				}
			}
		}
	}
	return p, nil
}

// Pull pulls the latest changes from the remote for a given repo.
// If the repo is not cloned, it will be cloned.
func (w Workdir) Pull(org, repo string) error {
	// Check if the repo is cloned.
	repoPath := filepath.Join(string(w), org, repo)
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		// Clone the repo.
		err = w.CloneGithubRepo(org, repo)
		if err != nil {
			return fmt.Errorf("CloneGithubRepo: %s", err)
		}
	}
	cmd := exec.Command("git", "pull")
	cmd.Dir = filepath.Dir(repoPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull: %s: %s", err, out)
	}
	return nil
}
