package shellquote_test

import (
	"fmt"

	"github.com/azghr/forge/shellquote"
)

func ExampleQuote() {
	unsafe := "file; rm -rf /"
	fmt.Println(shellquote.Quote(unsafe))
	// Output: 'file; rm -rf /'
}

func ExampleQuoteCommand() {
	cmd := shellquote.QuoteCommand([]string{"ls", "-l", "my file"})
	fmt.Println(cmd)
	// Output: 'ls' '-l' 'my file'
}

func ExampleQuoteWindows() {
	fmt.Println(shellquote.QuoteWindows(`C:\path with spaces\`))
	// Output: "C:\path with spaces\\"
}
