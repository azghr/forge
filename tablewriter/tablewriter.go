package tablewriter

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// Alignment specifies how a column's values are aligned.
type Alignment int

const (
	AlignLeft   Alignment = iota // left-aligned (default)
	AlignRight                   // right-aligned
	AlignCenter                  // center-aligned
)

// Table holds column definitions and rows for rendering tabular data.
type Table struct {
	mu      sync.RWMutex
	headers []string
	rows    [][]string
	padding int
	align   []Alignment
}

// New creates a new table with the given column headers.
func New(headers ...string) *Table {
	return NewWithOptions(headers)
}

// NewWithOptions creates a new table with headers and functional options.
func NewWithOptions(headers []string, opts ...Option) *Table {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	align := cfg.align
	if len(align) == 0 {
		align = make([]Alignment, len(headers))
	}

	return &Table{
		headers: headers,
		padding: cfg.padding,
		align:   align,
	}
}

// Append adds a row of values. It panics if the row length differs from
// the number of columns defined by headers.
func (t *Table) Append(row ...string) {
	if len(row) != len(t.headers) {
		panic(fmt.Sprintf("tablewriter: row has %d columns, expected %d", len(row), len(t.headers)))
	}
	r := make([]string, len(row))
	copy(r, row)

	t.mu.Lock()
	defer t.mu.Unlock()
	t.rows = append(t.rows, r)
}

// Render returns the table as an ASCII string.
//
// Each column cell is padded to max_content_width + 2*padding. Cells are
// separated by | in data rows and + in the separator line.
func (t *Table) Render() string {
	if len(t.headers) == 0 {
		return ""
	}

	t.mu.RLock()
	widths := t.columnWidths()
	var b strings.Builder

	t.writeRow(&b, widths, t.headers)
	t.writeSeparator(&b, widths)
	for _, row := range t.rows {
		t.writeRow(&b, widths, row)
	}
	t.mu.RUnlock()

	return trimTrailingSpaces(b.String())
}

func trimTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}
	return strings.Join(lines, "\n")
}

// Write writes the rendered table to w. Returns the number of bytes written.
func (t *Table) Write(w io.Writer) (int64, error) {
	s := t.Render()
	n, err := io.WriteString(w, s)
	return int64(n), err
}

// Len returns the number of data rows (excluding the header).
func (t *Table) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.rows)
}

func (t *Table) columnWidths() []int {
	widths := make([]int, len(t.headers))
	for i, h := range t.headers {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, v := range row {
			if len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}
	return widths
}

func (t *Table) writeRow(b *strings.Builder, widths []int, row []string) {
	for i, v := range row {
		align := AlignLeft
		if i < len(t.align) {
			align = t.align[i]
		}
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(alignCell(v, widths[i], align, t.padding))
	}
	b.WriteByte('\n')
}

func (t *Table) writeSeparator(b *strings.Builder, widths []int) {
	for i, w := range widths {
		if i > 0 {
			b.WriteByte('+')
		}
		total := w + 2*t.padding
		b.WriteString(strings.Repeat("-", total))
	}
	b.WriteByte('\n')
}

func alignCell(v string, width int, align Alignment, padding int) string {
	total := width + 2*padding
	switch align {
	case AlignRight:
		return fmt.Sprintf("%*s", total, v)
	case AlignCenter:
		left := padding + (width-len(v))/2
		right := total - left - len(v)
		if right < 0 {
			right = 0
		}
		return strings.Repeat(" ", left) + v + strings.Repeat(" ", right)
	default: // AlignLeft
		return strings.Repeat(" ", padding) + fmt.Sprintf("%-*s", width+padding, v)
	}
}
