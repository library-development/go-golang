package golang

import "github.com/library-development/go-english"

type Name struct {
	English    english.Name
	IsExported bool
}

func (n Name) String() string {
	if n.IsExported {
		return n.English.PascalCase()
	} else {
		return n.English.CamelCase()
	}
}
