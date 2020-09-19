package haystack

import "strings"

type Grid struct {
	meta Dict
	cols []Col
	rows []Row
}

// Meta returns the grid-level metadata
func (grid *Grid) Meta() Dict {
	return grid.meta
}

// Cols returns the column objects
func (grid *Grid) Cols() []Col {
	return grid.cols
}

// Rows returns the row objects
func (grid *Grid) Rows() []Row {
	return grid.rows
}

// ToZinc representes the object as:
//     ver:"3.0" <meta>
//     <col1>, <col2>, ...
//     <row1>
//     <row2>
//     ...
func (grid Grid) ToZinc() string {
	buf := strings.Builder{}
	grid.encodeTo(&buf, 0)
	return buf.String()
}

// Format as:
//     ver:"3.0" <meta>
//     <col1>, <col2>, ...
//     <row1>
//     <row2>
//     ...
//
// indentSize is the number of spaces to add to each new-line.
func (grid *Grid) encodeTo(buf *strings.Builder, indentSize int) {
	indentBuf := strings.Builder{}
	for i := 0; i < indentSize; i++ {
		indentBuf.WriteString(" ")
	}
	indent := indentBuf.String()

	buf.WriteString(indent)
	buf.WriteString("ver:\"3.0\"")
	if !grid.meta.isEmpty() {
		buf.WriteString(" ")
		grid.meta.encodeTo(buf, false)
	}
	buf.WriteString("\n")
	buf.WriteString(indent)
	for colIdx, col := range grid.cols {
		if colIdx != 0 {
			buf.WriteString(", ")
		}
		col.encodeTo(buf)
	}
	buf.WriteString("\n")
	buf.WriteString(indent)
	for rowIdx, row := range grid.rows {
		if rowIdx != 0 {
			buf.WriteString("\n")
			buf.WriteString(indent)
		}
		row.encodeTo(buf, &grid.cols)
	}
}

type Col struct {
	index int
	name  string
	meta  Dict
}

// Name returns the string name
func (col *Col) Name() string {
	return col.name
}

// Meta returns the column-level metadata
func (col *Col) Meta() Dict {
	return col.meta
}

// Format as <name> <meta>
func (col *Col) encodeTo(buf *strings.Builder) {
	buf.WriteString(col.name + " ")
	col.meta.encodeTo(buf, false)
}

type Row struct {
	items map[string]Val
}

// ToDict returns the values in a Dict format
func (row *Row) ToDict() Dict {
	return Dict{items: row.items}
}

// Format as <val1>, <val2>, ... Cols sets ordering
func (row *Row) encodeTo(buf *strings.Builder, cols *[]Col) {
	for colIdx, col := range *cols {
		if colIdx != 0 {
			buf.WriteString(", ")
		}
		item := row.items[col.name]
		itemZinc := item.ToZinc()
		buf.WriteString(itemZinc)
	}
}
