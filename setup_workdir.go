package golang

import (
	"fmt"
	"os"
	"os/exec"
)

func SetupWorkdir(loc string) (Workdir, error) {
	err := os.MkdirAll(loc, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("MkdirAll: %s", err)
	}
	cmd := exec.Command("go", "work", "init")
	cmd.Dir = loc
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go work init: %s: %s", err, out)
	}
	return Workdir(loc), nil
}
