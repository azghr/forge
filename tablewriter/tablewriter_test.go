package tablewriter_test

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/azghr/forge/tablewriter"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"Name", "Age"})
	if tbl.Len() != 0 {
		t.Errorf("expected 0 rows, got %d", tbl.Len())
	}
}

func TestAppendAndRender(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"Name", "Age"})
	tbl.Append("Alice", "30")
	tbl.Append("Bob", "25")

	out := tbl.Render()
	if !strings.Contains(out, "Alice") {
		t.Error("output missing Alice")
	}
	if !strings.Contains(out, "30") {
		t.Error("output missing 30")
	}
}

func TestRenderFormat(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"Name", "Age"})
	tbl.Append("Alice", "30")
	tbl.Append("Bob", "25")

	out := tbl.Render()
	const expected = " Name  | Age\n-------+-----\n Alice | 30\n Bob   | 25\n"
	if out != expected {
		t.Errorf("unexpected output:\n%q\nwant:\n%q", out, expected)
	}
}

func TestEmptyHeaders(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New(nil)
	out := tbl.Render()
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestNoRows(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"A", "B"})
	out := tbl.Render()
	if !strings.Contains(out, "A") || !strings.Contains(out, "B") {
		t.Error("headers should appear even with no rows")
	}
}

func TestAppendPanicOnMismatch(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"A", "B"})

	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		tbl.Append("1", "2", "3")
	}()
	if !didPanic {
		t.Error("expected panic for mismatched columns")
	}
}

func TestSingleColumn(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"X"})
	tbl.Append("a")
	tbl.Append("b")

	out := tbl.Render()
	if !strings.Contains(out, "X") || !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Error("single column table missing values")
	}
}

func TestSingleRow(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"A", "B"})
	tbl.Append("x", "y")

	out := tbl.Render()
	const expected = " A | B\n---+---\n x | y\n"
	if out != expected {
		t.Errorf("unexpected output:\n%q\nwant:\n%q", out, expected)
	}
}

func TestLongValues(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"Short", "VeryLongColumnName"})
	tbl.Append("a", "b")

	out := tbl.Render()
	if !strings.Contains(out, "VeryLongColumnName") {
		t.Error("long header missing")
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New([]string{"A"})
	tbl.Append("1")

	var buf bytes.Buffer
	n, err := tbl.Write(&buf)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n == 0 {
		t.Error("expected bytes written")
	}
	if buf.String() != tbl.Render() {
		t.Error("Write output differs from Render")
	}
}

func TestConcurrency(t *testing.T) {
	tbl := tablewriter.New([]string{"A", "B"})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tbl.Append("x", "y")
		}()
	}
	wg.Wait()

	if tbl.Len() != 10 {
		t.Errorf("expected 10 rows, got %d", tbl.Len())
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	tbl := tablewriter.New([]string{"A", "B"})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tbl.Append("x", "y")
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tbl.Render()
		}()
	}
	wg.Wait()
}

func TestTableDriven(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		headers  []string
		rows     [][]string
		opts     []tablewriter.Option
		wantSub  []string
		wantFull string
	}{
		{
			name:    "basic two columns",
			headers: []string{"X", "Y"},
			rows:    [][]string{{"foo", "bar"}},
			wantSub: []string{"X", "Y", "foo", "bar"},
		},
		{
			name:    "multiple rows",
			headers: []string{"A"},
			rows:    [][]string{{"1"}, {"2"}, {"3"}},
			wantSub: []string{"A", "1", "2", "3"},
		},
		{
			name:    "no rows",
			headers: []string{"H"},
			rows:    nil,
			wantSub: []string{"H"},
		},
		{
			name:     "with padding option",
			headers:  []string{"A", "B"},
			rows:     [][]string{{"x", "y"}},
			opts:     []tablewriter.Option{tablewriter.WithPadding(2)},
			wantFull: "  A  |  B\n-----+-----\n  x  |  y\n",
		},
		{
			name:     "right alignment",
			headers:  []string{"A", "B"},
			rows:     [][]string{{"x", "y"}},
			opts:     []tablewriter.Option{tablewriter.WithAlignment(tablewriter.AlignRight, tablewriter.AlignRight)},
			wantFull: "  A|  B\n---+---\n  x|  y\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := tablewriter.New(tt.headers, tt.opts...)
			for _, row := range tt.rows {
				tbl.Append(row...)
			}
			out := tbl.Render()

			if tt.wantFull != "" {
				if out != tt.wantFull {
					t.Errorf("unexpected output:\n%q\nwant:\n%q", out, tt.wantFull)
				}
			}
			for _, sub := range tt.wantSub {
				if !strings.Contains(out, sub) {
					t.Errorf("output missing %q", sub)
				}
			}
		})
	}
}

func TestWithAlignmentExact(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New(
		[]string{"Name", "Age", "City"},
		tablewriter.WithAlignment(tablewriter.AlignLeft, tablewriter.AlignRight, tablewriter.AlignCenter),
	)
	tbl.Append("Alice", "30", "NYC")
	tbl.Append("Bob", "25", "LA")

	_ = tbl.Render()
}

func TestWithAlignmentPartial(t *testing.T) {
	t.Parallel()

	tbl := tablewriter.New(
		[]string{"A", "B", "C"},
		tablewriter.WithAlignment(tablewriter.AlignRight),
	)
	tbl.Append("1", "2", "3")

	out := tbl.Render()
	if !strings.Contains(out, "1") || !strings.Contains(out, "2") || !strings.Contains(out, "3") {
		t.Error("partial alignment missing values")
	}
}

func BenchmarkRender(b *testing.B) {
	tbl := tablewriter.New([]string{"Name", "Age", "City"})
	for range 100 {
		tbl.Append("Alice", "30", "New York")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tbl.Render()
	}
}

func BenchmarkAppend(b *testing.B) {
	tbl := tablewriter.New([]string{"A", "B", "C"})
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		tbl.Append("x", "y", "z")
	}
}
