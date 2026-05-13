package tablewriter

// Option configures a Table.
type Option func(*config)

type config struct {
	padding int
	align   []Alignment
}

func defaultConfig() config {
	return config{
		padding: 1,
	}
}

// WithPadding sets the number of spaces on each side of a cell's content.
// The default is 1.
func WithPadding(n int) Option {
	return func(c *config) {
		c.padding = n
	}
}

// WithAlignment sets the per-column alignment. Pass one Alignment per column.
// If fewer alignments are provided than columns, remaining columns use
// AlignLeft.
//
// Example:
//
//	t := tablewriter.New(
//	    []string{"Name", "Age"},
//	    tablewriter.WithAlignment(tablewriter.AlignRight, tablewriter.AlignLeft),
//	)
func WithAlignment(align ...Alignment) Option {
	return func(c *config) {
		c.align = align
	}
}
