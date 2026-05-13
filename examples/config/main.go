package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/azghr/forge/atomicfile"
	"github.com/azghr/forge/jsonmerge"
	"github.com/azghr/forge/pathsafe"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: config <base.json> <override.json> [output.json]")
	}
	basePath := os.Args[1]
	overridePath := os.Args[2]
	outputPath := "merged.json"
	if len(os.Args) > 3 {
		outputPath = os.Args[3]
	}

	safeOut, err := pathsafe.SafeJoin(".", outputPath)
	if err != nil {
		log.Fatalf("unsafe output path: %v", err)
	}

	base := readJSON(basePath)
	override := readJSON(overridePath)

	jsonmerge.Merge(base, override)

	data, err := json.MarshalIndent(base, "", "  ")
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}

	if err := atomicfile.WriteFile(safeOut, bytes.NewReader(data)); err != nil {
		log.Fatalf("write: %v", err)
	}
	fmt.Printf("Merged config written to %s\n", safeOut)
}

func readJSON(path string) map[string]interface{} {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %q: %v", path, err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("parse %q: %v", path, err)
	}
	return m
}
