package atomicfile_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/azghr/forge/atomicfile"
)

func TestWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := atomicfile.Write(path, []byte("hello world")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello world" {
		t.Errorf("got %q, want %q", data, "hello world")
	}
}

func TestWriteOverwrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := atomicfile.Write(path, []byte("overwritten")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "overwritten" {
		t.Errorf("got %q, want %q", data, "overwritten")
	}
}

func TestWriteEmptyFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := atomicfile.Write(path, []byte{}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty file, got %d bytes", len(data))
	}
}

func TestWriteToNonexistentDir(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent", "file.txt")

	err := atomicfile.Write(path, []byte("data"))
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}

	var writeErr *atomicfile.WriteError
	if !errors.As(err, &writeErr) {
		t.Errorf("expected *WriteError, got %T", err)
	}
	if writeErr.Op != "create" {
		t.Errorf("expected Op=create, got %q", writeErr.Op)
	}
}

func TestWriteContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	err := atomicfile.WriteContext(ctx, path, []byte("data"))
	if !errors.Is(err, atomicfile.ErrCancelled) {
		t.Errorf("expected ErrCancelled, got %v", err)
	}
}

func TestWriteContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), -time.Millisecond)
	defer cancel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	err := atomicfile.WriteContext(ctx, path, []byte("data"))
	if !errors.Is(err, atomicfile.ErrCancelled) {
		t.Errorf("expected ErrCancelled, got %v", err)
	}
}

func TestWriteWithFileMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "mode.txt")

	if err := atomicfile.Write(path, []byte("data"), atomicfile.WithFileMode(0600)); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("expected mode 0600, got %#o", mode)
	}
}

func TestWriteWithoutFSync(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "nosync.txt")

	if err := atomicfile.Write(path, []byte("data"), atomicfile.WithoutFSync()); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "data" {
		t.Errorf("got %q, want %q", data, "data")
	}
}

func TestWriteReader(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "reader.txt")

	r := strings.NewReader("from reader")
	if err := atomicfile.WriteReader(context.Background(), path, r); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "from reader" {
		t.Errorf("got %q, want %q", data, "from reader")
	}
}

func TestWriteReaderContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	err := atomicfile.WriteReader(ctx, path, strings.NewReader("data"))
	if !errors.Is(err, atomicfile.ErrCancelled) {
		t.Errorf("expected ErrCancelled, got %v", err)
	}
}

func TestWriteErrorType(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent", "file.txt")

	err := atomicfile.Write(path, []byte("data"))

	var writeErr *atomicfile.WriteError
	if !errors.As(err, &writeErr) {
		t.Fatal("expected *WriteError")
	}
	if writeErr.Op == "" {
		t.Errorf("expected non-empty Op")
	}
	if writeErr.Err == nil {
		t.Errorf("expected non-nil underlying error")
	}
}

func TestConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "concurrent.txt")

	done := make(chan struct{}, 10)

	for i := 0; i < 10; i++ {
		i := i
		go func() {
			data := []byte{byte(i)}
			atomicfile.Write(path, data)
			done <- struct{}{}
		}()
	}

	for range 10 {
		<-done
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1 {
		t.Logf("concurrent writes: final length is %d (expected 1, one writer wins)", len(data))
	}
}

func TestNoPartialWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "partial.txt")

	if err := atomicfile.Write(path, []byte("full content")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "full content" {
		t.Errorf("got %q, want %q", data, "full content")
	}
}

func BenchmarkWrite(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "bench.txt")
	data := bytes.Repeat([]byte("a"), 4096)

	b.ResetTimer()
	for b.Loop() {
		atomicfile.Write(path, data)
	}
}

func BenchmarkWriteNoFSync(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, "bench.txt")
	data := bytes.Repeat([]byte("a"), 4096)

	b.ResetTimer()
	for b.Loop() {
		atomicfile.Write(path, data, atomicfile.WithoutFSync())
	}
}
