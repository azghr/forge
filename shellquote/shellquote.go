// Package shellquote escapes strings for safe use as shell arguments on POSIX
// or Windows. It wraps known logic to avoid injection issues.
package shellquote

import (
	"strings"
)

// Quote returns a shell-escaped version of s for POSIX shells.
// The string is wrapped in single quotes; any single quote inside the string
// is escaped by terminating the quote, adding an escaped quote, and resuming
// the quote (the '"'"' pattern).
func Quote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

// QuoteCommand escapes each argument in args and joins them into a single
// command-line string suitable for POSIX shells. Each argument is quoted
// via Quote; arguments are separated by a single space.
func QuoteCommand(args []string) string {
	if len(args) == 0 {
		return ""
	}
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(Quote(arg))
	}
	return b.String()
}

// QuoteWindows escapes s for Windows cmd.exe using double-quote wrapping and
// backslash escaping.
//
// The string is wrapped in double quotes. Backslashes preceding a double quote
// or at the end of the string are doubled; double quotes are escaped with a
// preceding backslash. This produces a string that CommandLineToArgv (and thus
// cmd.exe) will parse back into the original value.
func QuoteWindows(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for i := 0; i < len(s); {
		c := s[i]
		if c == '\\' {
			j := i
			for j < len(s) && s[j] == '\\' {
				j++
			}
			n := j - i
			if j == len(s) || s[j] == '"' {
				for k := 0; k < n*2; k++ {
					b.WriteByte('\\')
				}
			} else {
				for k := 0; k < n; k++ {
					b.WriteByte('\\')
				}
			}
			i = j
		} else if c == '"' {
			b.WriteString("\\\"")
			i++
		} else {
			b.WriteByte(c)
			i++
		}
	}
	b.WriteByte('"')
	return b.String()
}
