package golang

import "path/filepath"

type Ident struct {
	From string
	Name string
}

func (i Ident) String() string {
	return filepath.Join(i.From, i.Name)
}
