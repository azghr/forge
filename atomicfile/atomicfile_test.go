package atomicfile_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/azghr/forge/atomicfile"
)

type faultyReader struct{}

func (faultyReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("fail")
}

func TestWriteReplace(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fname := filepath.Join(dir, "test.txt")

	// Test successful write
	if err := atomicfile.WriteFile(fname, strings.NewReader("hello")); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(fname)
	if string(data) != "hello" {
		t.Errorf("Content = %s", data)
	}

	// Test overwrite existing
	if err := atomicfile.WriteFile(fname, strings.NewReader("world")); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(fname)
	if string(data) != "world" {
		t.Errorf("Overwrite failed")
	}

	// Test replace
	src := filepath.Join(dir, "src.txt")
	if err := os.WriteFile(src, []byte("xyz"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := atomicfile.ReplaceFile(src, fname); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(fname)
	if string(data) != "xyz" {
		t.Errorf("Replace failed")
	}
}

func TestAtomicity(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fname := filepath.Join(dir, "test.txt")

	// Write initial content
	if err := os.WriteFile(fname, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	// Faulty reader should cause WriteFile to fail
	bad := &faultyReader{}
	err := atomicfile.WriteFile(fname, bad)
	if err == nil {
		t.Fatal("expected error")
	}

	// Ensure original file unchanged
	data, _ := os.ReadFile(fname)
	if string(data) != "old" {
		t.Errorf("Expected old content, got %s", data)
	}
}

func TestWriteFileOptions(t *testing.T) {
	t.Parallel()

	t.Run("file mode", func(t *testing.T) {
		dir := t.TempDir()
		fname := filepath.Join(dir, "mode.txt")

		if err := atomicfile.WriteFile(fname, strings.NewReader("data"), atomicfile.WithFileMode(0600)); err != nil {
			t.Fatal(err)
		}

		info, err := os.Stat(fname)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0600 {
			t.Errorf("expected 0600, got %#o", info.Mode().Perm())
		}
	})

	t.Run("without fsync", func(t *testing.T) {
		dir := t.TempDir()
		fname := filepath.Join(dir, "nosync.txt")

		if err := atomicfile.WriteFile(fname, strings.NewReader("data"), atomicfile.WithoutFSync()); err != nil {
			t.Fatal(err)
		}

		data, _ := os.ReadFile(fname)
		if string(data) != "data" {
			t.Errorf("got %q", data)
		}
	})
}

func TestWriteFileToNonexistentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fname := filepath.Join(dir, "nonexistent", "test.txt")

	err := atomicfile.WriteFile(fname, strings.NewReader("data"))
	if err == nil {
		t.Fatal("expected error")
	}

	var writeErr *atomicfile.WriteError
	if !errors.As(err, &writeErr) {
		t.Errorf("expected *WriteError, got %T", err)
	}
	if writeErr.Op != "create" {
		t.Errorf("expected Op=create, got %q", writeErr.Op)
	}
}

func TestReplaceFileError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fname := filepath.Join(dir, "test.txt")
	src := filepath.Join(dir, "nonexistent-src")

	err := atomicfile.ReplaceFile(src, fname)
	if err == nil {
		t.Fatal("expected error for nonexistent source")
	}

	var writeErr *atomicfile.WriteError
	if !errors.As(err, &writeErr) {
		t.Errorf("expected *WriteError, got %T", err)
	}
	if writeErr.Op != "rename" {
		t.Errorf("expected Op=rename, got %q", writeErr.Op)
	}
}

func TestConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	fname := filepath.Join(dir, "concurrent.txt")

	done := make(chan struct{}, 10)

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			data := fmt.Sprintf("data-%d", i)
			atomicfile.WriteFile(fname, strings.NewReader(data))
			done <- struct{}{}
		}()
	}

	for range 10 {
		<-done
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("expected some content after concurrent writes")
	}
}

func TestWriteFileEmpty(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	fname := filepath.Join(dir, "empty.txt")

	if err := atomicfile.WriteFile(fname, strings.NewReader("")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(data))
	}
}

func BenchmarkWriteFile(b *testing.B) {
	dir := b.TempDir()
	fname := filepath.Join(dir, "bench.txt")
	r := bytes.NewReader(bytes.Repeat([]byte("a"), 4096))

	b.ResetTimer()
	for b.Loop() {
		r.Reset(bytes.Repeat([]byte("a"), 4096))
		atomicfile.WriteFile(fname, r)
	}
}

func BenchmarkWriteFileNoFSync(b *testing.B) {
	dir := b.TempDir()
	fname := filepath.Join(dir, "bench.txt")
	r := bytes.NewReader(bytes.Repeat([]byte("a"), 4096))

	b.ResetTimer()
	for b.Loop() {
		r.Reset(bytes.Repeat([]byte("a"), 4096))
		atomicfile.WriteFile(fname, r, atomicfile.WithoutFSync())
	}
}
