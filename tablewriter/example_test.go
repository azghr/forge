package tablewriter_test

import (
	"fmt"

	"github.com/azghr/forge/tablewriter"
)

func Example() {
	tbl := tablewriter.New("Name", "Age")
	tbl.Append("Alice", "30")
	tbl.Append("Bob", "25")
	fmt.Print(tbl.Render())
	// Output:
	//  Name  | Age
	// -------+-----
	//  Alice | 30
	//  Bob   | 25
}

func Example_singleColumn() {
	tbl := tablewriter.New("Score")
	tbl.Append("42")
	tbl.Append("100")
	fmt.Print(tbl.Render())
	// Output:
	//  Score
	// -------
	//  42
	//  100
}

func Example_padding() {
	tbl := tablewriter.NewWithOptions(
		[]string{"A", "B"},
		tablewriter.WithPadding(2),
	)
	tbl.Append("x", "y")
	fmt.Print(tbl.Render())
	// Output:
	//   A  |  B
	// -----+-----
	//   x  |  y
}

func Example_alignment() {
	tbl := tablewriter.NewWithOptions(
		[]string{"Item", "Price"},
		tablewriter.WithAlignment(tablewriter.AlignLeft, tablewriter.AlignRight),
	)
	tbl.Append("Apple", "$1.50")
	tbl.Append("Banana", "$0.75")
	fmt.Print(tbl.Render())
	// Output:
	//  Item   |  Price
	// --------+-------
	//  Apple  |  $1.50
	//  Banana |  $0.75
}
