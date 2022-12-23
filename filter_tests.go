package golang

import (
	"io/fs"
	"strings"
)

func FilterTests(fi fs.FileInfo) bool {
	if strings.HasSuffix(fi.Name(), "_test.go") {
		return false
	}
	return true
}
