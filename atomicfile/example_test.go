package atomicfile_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/azghr/forge/atomicfile"
)

func ExampleWriteFile() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "config.txt")
	data := bytes.NewBufferString("important")

	if err := atomicfile.WriteFile(fname, data); err != nil {
		fmt.Println("write failed:", err)
		return
	}
	// config.txt is fully written or untouched.

	content, _ := os.ReadFile(fname)
	fmt.Println(string(content))
	// Output: important
}

func ExampleReplaceFile() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	os.WriteFile(src, []byte("hello"), 0644)
	os.WriteFile(dst, []byte("old"), 0644)

	if err := atomicfile.ReplaceFile(src, dst); err != nil {
		fmt.Println("replace failed:", err)
		return
	}

	content, _ := os.ReadFile(dst)
	fmt.Println(string(content))
	// Output: hello
}

func ExampleWithFileMode() {
	dir, _ := os.MkdirTemp("", "atomicfile-example")
	defer os.RemoveAll(dir)

	fname := filepath.Join(dir, "secure.txt")

	if err := atomicfile.WriteFile(fname, strings.NewReader("data"), atomicfile.WithFileMode(0600)); err != nil {
		fmt.Println("write failed:", err)
		return
	}

	info, _ := os.Stat(fname)
	fmt.Printf("%#o\n", info.Mode().Perm())
	// Output: 0600
}
