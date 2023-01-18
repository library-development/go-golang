package golang

import (
	"github.com/library-development/go-english"
	"github.com/library-development/go-nameconv"
)

func ParseName(name string) (english.Name, error) {
	n, err := nameconv.ParseCamelCase(name)
	if err != nil {
		n, err = nameconv.ParsePascalCase(name)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}
