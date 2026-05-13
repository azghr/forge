package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScaffoldExample(t *testing.T) {
	tmpDir := t.TempDir()
	err := scaffoldExample("cli", tmpDir)
	if err != nil {
		t.Fatalf("scaffoldExample(cli): %v", err)
	}

	expectedFiles := []string{"main.go"}
	for _, f := range expectedFiles {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", f)
		}
	}
}

func TestScaffoldExampleServer(t *testing.T) {
	tmpDir := t.TempDir()
	err := scaffoldExample("server", tmpDir)
	if err != nil {
		t.Fatalf("scaffoldExample(server): %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "main.go")); os.IsNotExist(err) {
		t.Error("expected main.go to exist")
	}
}

func TestScaffoldExample_invalidType(t *testing.T) {
	err := scaffoldExample("bogus", t.TempDir())
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !strings.Contains(err.Error(), "unknown example type") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestScaffoldNewCLI(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "mycli")
	err := scaffoldNew("cli", tmpDir)
	if err != nil {
		t.Fatalf("scaffoldNew(cli): %v", err)
	}

	expected := []string{"main.go", "go.mod", "Makefile"}
	for _, f := range expected {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", f)
		}
	}
}

func TestScaffoldNewServer(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "myserver")
	err := scaffoldNew("server", tmpDir)
	if err != nil {
		t.Fatalf("scaffoldNew(server): %v", err)
	}

	expected := []string{"main.go", "go.mod", "Dockerfile", "Makefile"}
	for _, f := range expected {
		path := filepath.Join(tmpDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", f)
		}
	}
}

func TestColor_noColorEnv(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	result := green("hello")
	if strings.Contains(result, "\033") {
		t.Errorf("expected no ANSI codes when NO_COLOR is set, got: %q", result)
	}
}

func TestColor_ansiCodes(t *testing.T) {
	os.Unsetenv("NO_COLOR")
	result := green("hello")
	if !strings.Contains(result, "\033[32m") {
		t.Errorf("expected ANSI green code, got: %q", result)
	}
}

func TestFindForgeRoot(t *testing.T) {
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	os.Chdir(t.TempDir())
	root := findForgeRoot()
	if root == "" {
		t.Error("expected non-empty root")
	}
	// When no go.work is found, should return "."
	if root != "." {
		t.Errorf("expected '.', got %q", root)
	}
}

func TestListExampleTypes(t *testing.T) {
	got := captureStdout(listExampleTypes)
	if !strings.Contains(got, "cli") || !strings.Contains(got, "server") {
		t.Errorf("expected cli and server in output, got: %s", got)
	}
}

func TestListNewTypes(t *testing.T) {
	got := captureStdout(listNewTypes)
	if !strings.Contains(got, "cli") || !strings.Contains(got, "server") {
		t.Errorf("expected cli and server in output, got: %s", got)
	}
}

func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	b, _ := io.ReadAll(r)
	os.Stdout = stdout
	return string(b)
}
