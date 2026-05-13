package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed template/example/*/*.tmpl
var exampleTemplates embed.FS

type exampleType struct {
	Name        string
	Description string
}

var exampleTypes = []exampleType{
	{Name: "cli", Description: "CLI tool using flagsub and envconfig"},
	{Name: "server", Description: "HTTP server with graceful shutdown"},
}

func runExample() {
	args := os.Args[2:]
	if len(args) == 0 {
		listExampleTypes()
		return
	}

	typeName := args[0]
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	if err := scaffoldExample(typeName, dir); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", red("error:"), err)
		os.Exit(1)
	}
}

func listExampleTypes() {
	fmt.Println(cyan("Available example types:"))
	for _, t := range exampleTypes {
		fmt.Printf("  %-12s %s\n", t.Name, t.Description)
	}
	fmt.Println()
	fmt.Printf("Usage: %s example <type> [output-dir]\n", os.Args[0])
}

func scaffoldExample(typeName, outDir string) error {
	valid := false
	for _, t := range exampleTypes {
		if t.Name == typeName {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("unknown example type %q", typeName)
	}

	tmplDir := fmt.Sprintf("template/example/%s", typeName)
	entries, err := fs.ReadDir(exampleTemplates, tmplDir)
	if err != nil {
		return fmt.Errorf("reading templates: %w", err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	data := map[string]string{
		"PackageName": filepath.Base(outDir),
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		tmplPath := filepath.Join(tmplDir, e.Name())
		tmplContent, err := exampleTemplates.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", e.Name(), err)
		}

		outName := strings.TrimSuffix(e.Name(), ".tmpl")
		outPath := filepath.Join(outDir, outName)

		tmpl, err := template.New(outName).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", e.Name(), err)
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("creating %s: %w", outName, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return fmt.Errorf("executing template %s: %w", e.Name(), err)
		}
		fmt.Printf("  %s %s\n", green("created"), outPath)
	}

	return nil
}
