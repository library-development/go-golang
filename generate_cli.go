package golang

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
)

// GenerateCLI geneates a simple Go CLI app that only calls a single func.
// The resulting command has no arguments.
// It reads function input from stdin in JSON format and writes the result to stdout in JSON format.
func GenerateCLI(srcDir, pkg, funcName, outFile string) error {
	funcSignature, err := ReadFuncSignature(srcDir, pkg, funcName)
	if err != nil {
		return err
	}
	importMap := ImportMap{}
	importMap.AddPackage(pkg)
	importMap.AddPackage("encoding/json")
	importMap.AddPackage("os")
	for _, input := range funcSignature.Inputs {
		for _, pkg := range input.Type.Packages() {
			importMap.AddPackage(pkg)
		}
	}
	for _, output := range funcSignature.Outputs {
		for _, pkg := range output.Type.Packages() {
			importMap.AddPackage(pkg)
		}
	}
	err = os.MkdirAll(filepath.Dir(outFile), os.ModePerm)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	b.WriteString("package main\n\n")
	importMap.Write(&b)
	b.WriteString("\nfunc main() {\n")
	b.WriteString("\tvar input struct {\n")
	for _, input := range funcSignature.Inputs {
		name, err := ParseName(input.Name)
		if err != nil {
			return err
		}
		b.WriteString("\t\t")
		b.WriteString(name.PascalCase())
		b.WriteString(" ")
		input.Type.Write(&b)
		b.WriteString(" `json:\"")
		b.WriteString(name.SnakeCase())
		b.WriteString("\"`\n")
	}
	b.WriteString("\t}\n")
	b.WriteString("\tjson.NewDecoder(os.Stdin).Decode(&input)\n")
	b.WriteString("\t")
	if len(funcSignature.Outputs) > 0 {
		for i := range funcSignature.Outputs {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("out")
			b.WriteString(strconv.Itoa(i + 1))
		}
		b.WriteString(" := ")
	}
	p, ok := importMap.Resolve(pkg)
	if !ok {
		panic("package not found")
	}
	b.WriteString(p)
	b.WriteString(".")
	b.WriteString(funcName)
	b.WriteString("(")
	for i, input := range funcSignature.Inputs {
		name, err := ParseName(input.Name)
		if err != nil {
			return err
		}
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("input.")
		b.WriteString(name.PascalCase())
	}
	b.WriteString(")\n")
	b.WriteString("\tvar output struct {\n")
	for _, output := range funcSignature.Outputs {
		name, err := ParseName(output.Name)
		if err != nil {
			return err
		}
		b.WriteString("\t\t")
		b.WriteString(name.PascalCase())
		b.WriteString(" ")
		output.Type.Write(&b)
		b.WriteString(" `json:\"")
		b.WriteString(name.SnakeCase())
		b.WriteString("\"`\n")
	}
	b.WriteString("\t}\n")
	for i := range funcSignature.Outputs {
		b.WriteString("\toutput.Out")
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(" = out")
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString("\n")
	}
	b.WriteString("\tjson.NewEncoder(os.Stdout).Encode(output)\n")
	b.WriteString("}\n")
	err = os.WriteFile(outFile, b.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}
	err = RunGoimports(outFile)
	if err != nil {
		return err
	}
	return nil
}
