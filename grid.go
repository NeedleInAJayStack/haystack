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

// EmptyGrid returns an empty grid
func EmptyGrid() Grid {
	return Grid{}
}

// Meta returns the grid-level metadata
func (grid *Grid) Meta() *Dict {
	return &grid.meta
}

// ColCount returns the count of columns
func (grid *Grid) ColCount() int {
	return len(grid.cols)
}

// Col returns the column matching the name
func (grid *Grid) Col(name string) *Col {
	var colMatch *Col
	for i := 0; i < grid.ColCount(); i++ {
		col := grid.ColAt(i)
		if grid.cols[i].name == name {
			colMatch = col
			break
		}
	}
	if colMatch == nil {
		panic("Unknown column: " + name)
	}
	return colMatch
}

// ColAt returns the column at the index
func (grid *Grid) ColAt(index int) *Col {
	return &grid.cols[index]
}

// RowCount returns the count of rows
func (grid *Grid) RowCount() int {
	return len(grid.rows)
}

// RowAt returns the row at the index
func (grid *Grid) RowAt(index int) *Row {
	return &grid.rows[index]
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
	if len(grid.cols) == 0 { // Empty grids get just the word: empty
		buf.WriteString("empty\n")
	} else {
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

func (row *Row) Get(name string) Val {
	return row.items[name]
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
