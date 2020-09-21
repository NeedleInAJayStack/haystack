package haystack

import (
	"bufio"
	"strings"
)

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
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	grid.WriteZincTo(out, 0)
	out.Flush()
	return builder.String()
}

// Format as:
//     ver:"3.0" <meta>
//     <col1>, <col2>, ...
//     <row1>
//     <row2>
//     ...
//
// indentSize is the number of spaces to add to each new-line.
func (grid *Grid) WriteZincTo(buf *bufio.Writer, indentSize int) {
	writeIndent(buf, indentSize)
	buf.WriteString("ver:\"3.0\"")
	if !grid.meta.IsEmpty() {
		buf.WriteString(" ")
		grid.meta.WriteZincTo(buf, false)
	}
	buf.WriteString("\n")
	writeIndent(buf, indentSize)
	for colIdx, col := range grid.cols {
		if colIdx != 0 {
			buf.WriteString(", ")
		}
		col.WriteZincTo(buf)
	}
	buf.WriteString("\n")
	writeIndent(buf, indentSize)
	for rowIdx, row := range grid.rows {
		if rowIdx != 0 {
			buf.WriteString("\n")
			writeIndent(buf, indentSize)
		}
		row.WriteZincTo(buf, &grid.cols, indentSize)
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
func (col *Col) WriteZincTo(buf *bufio.Writer) {
	buf.WriteString(col.name)
	if !col.meta.IsEmpty() {
		buf.WriteRune(' ')
		col.meta.WriteZincTo(buf, false)
	}
}

type Row struct {
	items map[string]Val
}

// ToDict returns the values in a Dict format
func (row *Row) ToDict() Dict {
	return Dict{items: row.items}
}

// Format as <val1>, <val2>, ... Cols sets ordering
func (row *Row) WriteZincTo(buf *bufio.Writer, cols *[]Col, indentSize int) {
	for colIdx, col := range *cols {
		if colIdx != 0 {
			buf.WriteString(", ")
		}
		val := row.items[col.name]
		switch val := val.(type) {
		case Grid:
			indentSize = indentSize + 1
			buf.WriteString("<<\n")
			val.WriteZincTo(buf, indentSize)
			buf.WriteString("\n")
			writeIndent(buf, indentSize)
			buf.WriteString(">>")
			indentSize = indentSize - 1
		case List:
			val.WriteZincTo(buf)
		case Dict:
			val.WriteZincTo(buf, true)
		case Str:
			val.WriteZincTo(buf)
		case Uri:
			val.WriteZincTo(buf)
		default:
			buf.WriteString(val.ToZinc())
		}
	}
}

func writeIndent(buf *bufio.Writer, indentSize int) {
	for i := 0; i < indentSize; i++ {
		buf.WriteString("  ") // Each indent is 2 spaces
	}
}
