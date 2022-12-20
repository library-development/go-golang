package golang

import (
	"fmt"
	"io"
	"path/filepath"
)

// ImportMap maps package paths to local variable names.
type ImportMap map[string]string

// AddPackage adds the given package path to the import map.
// It will ensure there are no naming conflicts by appending as many underscores as necessary to the end of the variable name.
func (i ImportMap) AddPackage(pkg string) string {
	name, ok := i[pkg]
	if ok {
		return name
	}
	name = filepath.Base(pkg)
	for {
		conflict := false
		for _, v := range i {
			if name == v {
				conflict = true
			}
		}
		if !conflict {
			break
		}
		name += "_"
	}
	i[pkg] = name
	return name
}

// Resolve will return the associate variable name with an imported package.
func (i ImportMap) Resolve(pkg string) (string, bool) {
	name, ok := i[pkg]
	return name, ok
}

// Write will write the import map in Go to w.
func (i ImportMap) Write(w io.Writer) error {
	if len(i) == 0 {
		return nil
	}
	_, err := w.Write([]byte("import (\n"))
	if err != nil {
		return err
	}
	for p, v := range i {
		_, err = fmt.Fprintf(w, "\t%s \"%s\"\n", v, p)
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte(")\n"))
	if err != nil {
		return err
	}
	return nil
}
