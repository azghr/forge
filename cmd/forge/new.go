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

//go:embed template/new/*/*.tmpl
var newTemplates embed.FS

type newType struct {
	Name        string
	Description string
}

var newTypes = []newType{
	{Name: "cli", Description: "Production-ready CLI tool"},
	{Name: "server", Description: "Production-ready HTTP server"},
	{Name: "api", Description: "Production HTTP API with forge packages"},
}

func runNew() {
	args := os.Args[2:]
	if len(args) < 1 {
		listNewTypes()
		return
	}

	typeName := args[0]
	projectName := ""
	if len(args) > 1 {
		projectName = args[1]
	} else {
		projectName = typeName
	}

	valid := false
	for _, t := range newTypes {
		if t.Name == typeName {
			valid = true
			break
		}
	}
	if !valid {
		fmt.Fprintf(os.Stderr, "%s unknown project type %q\n", red("error:"), typeName)
		listNewTypes()
		os.Exit(1)
	}

	if err := scaffoldNew(typeName, projectName); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", red("error:"), err)
		os.Exit(1)
	}
}

func listNewTypes() {
	fmt.Println(cyan("Available project types:"))
	for _, t := range newTypes {
		fmt.Printf("  %-12s %s\n", t.Name, t.Description)
	}
	fmt.Println()
	fmt.Printf("Usage: %s new <type> [project-name]\n", os.Args[0])
}

func scaffoldNew(typeName, projectName string) error {
	outDir := projectName
	tmplDir := fmt.Sprintf("template/new/%s", typeName)
	entries, err := fs.ReadDir(newTemplates, tmplDir)
	if err != nil {
		return fmt.Errorf("reading templates: %w", err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	baseName := filepath.Base(projectName)
	data := map[string]string{
		"ProjectName": baseName,
		"PackageName": strings.ReplaceAll(baseName, "-", ""),
		"ModulePath":  baseName,
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		tmplPath := filepath.Join(tmplDir, e.Name())
		tmplContent, err := newTemplates.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", e.Name(), err)
		}

		subDirs := ""
		if strings.Contains(e.Name(), "/") {
			subDirs = filepath.Dir(e.Name())
			if err := os.MkdirAll(filepath.Join(outDir, subDirs), 0755); err != nil {
				return fmt.Errorf("creating subdir: %w", err)
			}
		}

		outName := strings.TrimSuffix(e.Name(), ".tmpl")
		outPath := filepath.Join(outDir, subDirs, outName)

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

	fmt.Printf("\n%s Project %q scaffolded in %s/\n", green("done"), baseName, outDir)
	fmt.Println()
	fmt.Println(cyan("Next steps:"))
	fmt.Printf("  cd %s\n", outDir)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run .\n")
	if typeName == "api" || typeName == "server" {
		fmt.Printf("  make run\n")
	}
	return nil
}
