# stringutil

Common string transformations: Title (capitalise each word), Slug (URL-safe), RemoveAccents, etc.

## Problem

Go's standard library provides basic string manipulation in `strings` and `unicode`,
but common higher-level transformations are missing or deprecated (`strings.Title`
is deprecated and doesn't handle Unicode properly). This package fills those gaps
with zero external dependencies.

## Quick start

```go
import "github.com/azghr/forge/stringutil"

fmt.Println(stringutil.Title("hello world"))        // "Hello World"
fmt.Println(stringutil.Slug("Go Lang Library"))      // "go-lang-library"
fmt.Println(stringutil.RemoveAccents("café"))        // "cafe"
```

## API

### Functions

- **`Title(s string) string`** — capitalise every word (letter sequence) in `s`.
  Unlike `strings.Title` (deprecated), correctly handles Unicode punctuation.
- **`Slug(s string, opts ...SlugOption) string`** — produce a URL-safe slug:
  lowercased, non-alphanumeric characters replaced by a separator, consecutive
  separators collapsed, leading/trailing trimmed.
- **`RemoveAccents(s string) string`** — strip diacritical marks from `s` using
  an internal mapping of common accented characters to their ASCII base forms.

### Options

- **`WithSeparator(sep string) SlugOption`** — separator between words (default `"-"`).
- **`WithMaxLength(n int) SlugOption`** — limit slug to `n` bytes, truncating at
  the nearest separator boundary. `0` means unlimited (default).

## Performance

| Function       | Time | Space    | Notes                                     |
|----------------|------|----------|-------------------------------------------|
| Title          | O(n) | O(n)     | Operates on runes in-place                |
| Slug           | O(n) | O(n)     | Single pass, builder-backed               |
| RemoveAccents  | O(n) | O(n)     | Map lookup per rune; zero allocs for pure ASCII |

All functions are concurrency-safe (no shared mutable state).
