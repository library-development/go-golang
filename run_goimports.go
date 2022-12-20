package golang

import "os/exec"

func RunGoimports(path string) error {
	return exec.Command("goimports", "-w", path).Run()
}
