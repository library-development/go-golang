package golang

import (
	"github.com/library-development/go-english"
	"github.com/library-development/go-nameconv"
)

func EnglishName(name string) english.Name {
	n, err := nameconv.ParsePascalCase(name)
	if err != nil {
		n, err := nameconv.ParseCamelCase(name)
		if err != nil {
			panic(err)
		}
		return n
	}
	return n
}
