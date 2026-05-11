package atomicfile_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/azghr/forge/atomicfile"
)

func ExampleWrite() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "example.txt")

	if err := atomicfile.Write(path, []byte("hello world")); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	data, _ := os.ReadFile(path)
	fmt.Println(string(data))
	// Output: hello world
}

func ExampleWriteContext() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "example.txt")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := atomicfile.WriteContext(ctx, path, []byte("hello")); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	data, _ := os.ReadFile(path)
	fmt.Println(string(data))
	// Output: hello
}

func ExampleWriteReader() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "example.txt")
	r := strings.NewReader("from reader")

	if err := atomicfile.WriteReader(context.Background(), path, r); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	data, _ := os.ReadFile(path)
	fmt.Println(string(data))
	// Output: from reader
}

func ExampleWithFileMode() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "example.txt")

	if err := atomicfile.Write(path, []byte("data"), atomicfile.WithFileMode(0600)); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	info, _ := os.Stat(path)
	fmt.Printf("%#o\n", info.Mode().Perm())
	// Output: 0600
}
