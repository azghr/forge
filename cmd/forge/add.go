package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type forgePkg struct {
	Name        string
	Description string
}

var forgePackages = []forgePkg{
	{Name: "atomicfile", Description: "Atomic file writes"},
	{Name: "cache", Description: "Generic in-memory TTL cache"},
	{Name: "envconfig", Description: "Environment config loading"},
	{Name: "flagsub", Description: "Subcommand support for flag package"},
	{Name: "jsonmerge", Description: "Recursive JSON merge and diff"},
	{Name: "lockutil", Description: "TryLock and context-aware Lock"},
	{Name: "mathutil", Description: "Math helpers (Clamp, Sign, Lerp)"},
	{Name: "multityperror", Description: "Aggregate multiple errors"},
	{Name: "option", Description: "Generic Option (Maybe) type"},
	{Name: "orderedset", Description: "Insertion-ordered set"},
	{Name: "pathsafe", Description: "Safe path joining"},
	{Name: "priorityqueue", Description: "Generic binary heap"},
	{Name: "queue", Description: "Generic FIFO queue"},
	{Name: "regexcache", Description: "Compiled regex cache"},
	{Name: "retry", Description: "Retry with backoff and jitter"},
	{Name: "shellquote", Description: "Shell-safe string quoting"},
	{Name: "sliceutil", Description: "Slice operations (Map, Filter, Reduce)"},
	{Name: "stopwatch", Description: "Benchmarking stopwatch"},
	{Name: "stringutil", Description: "String transformations"},
	{Name: "tablewriter", Description: "ASCII table formatting"},
	{Name: "validator", Description: "Struct field validation"},
	{Name: "workerpool", Description: "Fixed-size worker pool"},
}

func runAdd() {
	args := os.Args[2:]
	if len(args) == 0 {
		listPackages()
		return
	}

	name := args[0]
	for _, p := range forgePackages {
		if p.Name == name {
			fmt.Printf("%s Run this in your project:\n\n", green("→"))
			fmt.Printf("  go get github.com/azghr/forge/%s@latest\n\n", name)
			fmt.Printf("  Package: %s\n", p.Name)
			fmt.Printf("  Purpose: %s\n", p.Description)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "%s unknown package %q\n", red("error:"), name)
	listPackages()
	os.Exit(1)
}

func listPackages() {
	fmt.Println(cyan("Available forge packages:"))
	sort.Slice(forgePackages, func(i, j int) bool {
		return forgePackages[i].Name < forgePackages[j].Name
	})
	for _, p := range forgePackages {
		fmt.Printf("  %-20s %s\n", p.Name, p.Description)
	}
	fmt.Println()
	fmt.Printf("Usage: %s add <package>\n", os.Args[0])
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  forge add retry")
	fmt.Println("  forge add envconfig")
	fmt.Println("  forge add stopwatch")
	fmt.Println()
	fmt.Println("Or add multiple packages to go.mod at once:")
	fmt.Println("  forge add retry | grep 'go get' | sh")
}

// pkgModulePath returns the full module path for a forge package.
func pkgModulePath(name string) string {
	return "github.com/azghr/forge/" + name
}

// findPkg attempts to find a forge package by name (case-insensitive prefix).
func findPkg(name string) *forgePkg {
	name = strings.ToLower(name)
	for _, p := range forgePackages {
		if strings.EqualFold(p.Name, name) {
			return &p
		}
	}
	return nil
}
