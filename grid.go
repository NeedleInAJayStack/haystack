package haystack

import "strings"

type Col struct {
	index int
	name  string
	meta  Dict
}

// Format as <name> <meta>
func (col *Col) encodeTo(buf *strings.Builder) {
	buf.WriteString(col.name + " ")
	col.meta.encodeTo(buf, false)
}

type Row struct {
	vals []Val
}

// Format as <val1>, <val2>, ...
func (row *Row) encodeTo(buf *strings.Builder) {
	for idx, val := range row.vals {
		if idx != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(val.toZinc())
	}
}

type Grid struct {
	meta Dict
	cols []Col
	rows []Row
}

// Format as:
//     ver:"3.0" <meta>
//     <col1>, <col2>, ...
//     <row1>
//     <row2>
//     ...
func (grid *Grid) toZinc() string {
	buf := strings.Builder{}
	buf.WriteString("ver:\"3.0\"")
	if !grid.meta.isEmpty() {
		buf.WriteString(" ")
		grid.meta.encodeTo(&buf, false)
	}
	buf.WriteString("\n")
	for colIdx, col := range grid.cols {
		if colIdx != 0 {
			buf.WriteString(", ")
		}
		col.encodeTo(&buf)
	}
	buf.WriteString("\n")
	for rowIdx, row := range grid.rows {
		if rowIdx != 0 {
			buf.WriteString("\n")
		}
		row.encodeTo(&buf)
	}
	return buf.String()
}
