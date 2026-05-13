# tablewriter

Format tabular data as ASCII tables.

## Problem

Presenting rows and columns of data in text form (CLI output, log files,
code comments) requires manually aligning columns with padding, separators,
and borders. `tablewriter` handles this automatically.

Features:
- **ASCII table** rendering with `|` and `-` separators.
- **Per-column alignment** — left, right, or center.
- **Configurable padding** — number of spaces around cell content.
- **Thread-safe** concurrent reads and appends.

## Quick start

```go
import "github.com/azghr/forge/tablewriter"

tbl := tablewriter.New([]string{"Name", "Age"})
tbl.Append("Alice", "30")
tbl.Append("Bob", "25")
fmt.Println(tbl.Render())

// Name | Age
// -----+----
// Alice | 30
// Bob   | 25
```

## API

### Types

- **`Table`** — holds headers and rows.
- **`Alignment`** — `AlignLeft`, `AlignRight`, `AlignCenter`.

### Functions

- **`New(headers []string, opts ...Option) *Table`** — creates a table with
  headers and optional configuration (padding, alignment).

### Methods

- **`Append(row ...string)`** — adds a row. Panics on column count mismatch.
- **`Render() string`** — returns the formatted table as a string.
- **`Write(w io.Writer) (int64, error)`** — writes the rendered table.
- **`Len() int`** — number of data rows.

### Options

- **`WithPadding(n int)`** — set cell padding (default 1).
- **`WithAlignment(align ...Alignment)`** — set per-column alignment.

## Performance

- `Render` builds the full output string: O(rows × cols).
- `Append` copies the row slice: O(cols).
- No allocations per cell beyond the string building in `Render`.
