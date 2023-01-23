package golang_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/library-development/go-golang"
)

func TestParsePackage(t *testing.T) {
	dir, err := os.MkdirTemp("", "golang")
	if err != nil {
		t.Fatal(err)
	}
	w, err := golang.SetupWorkdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = w.CloneGithubRepo("library-development", "go-golang")
	if err != nil {
		t.Fatal(err)
	}
	pkg, err := w.ParsePackage("github.com/library-development/go-golang")
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(pkg)
}
