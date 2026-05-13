package pathsafe_test

import (
	"fmt"

	"github.com/azghr/forge/pathsafe"
)

func ExampleSafeJoin() {
	p, err := pathsafe.SafeJoin("/home/user", "docs/report.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p)
	// Output: /home/user/docs/report.pdf
}

func ExampleSafeJoin_traversal() {
	_, err := pathsafe.SafeJoin("/home/user", "../etc/passwd")
	fmt.Println(err)
	// Output: pathsafe: result path is outside base directory
}

func ExampleSafeJoin_symlinkOption() {
	// AllowSymlinkFollow resolves symlinks before the containment check.
	// This prevents symlink-based traversal. When the path does not exist,
	// an error is returned.
	_, err := pathsafe.SafeJoin(
		"/home/user",
		"docs/report.pdf",
		pathsafe.AllowSymlinkFollow(),
	)
	fmt.Println(err != nil)
	// Output: true
}
