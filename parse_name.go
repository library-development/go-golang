package golang

import "lib.dev/nameconv"

func ParseName(name string) (*nameconv.Name, error) {
	n, err := nameconv.ParseCamelCase(name)
	if err != nil {
		n, err = nameconv.ParsePascalCase(name)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
}
