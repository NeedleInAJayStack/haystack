package haystack

import (
	"bufio"
	"encoding/json"
	"strings"
)

// Grid is a two dimension data structure of cols and rows.
type Grid struct {
	meta Dict
	cols []Col
	rows []Row
}

// EmptyGrid creates an empty grid.
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

// MarshalJSON represents the object in a special JSON object format. See https://project-haystack.org/doc/Json#grid
func (grid Grid) MarshalJSON() ([]byte, error) {
	buf := strings.Builder{}

	buf.WriteString("{\"meta\":")
	newMeta := grid.meta.Set("ver", NewStr("3.0")) // Add in version
	metaJson, metaErr := json.Marshal(newMeta)
	if metaErr != nil {
		return []byte{}, metaErr
	}
	buf.WriteString(string(metaJson))

	buf.WriteString(",\"cols\":")
	colsJson, colsErr := json.Marshal(grid.cols)
	if colsErr != nil {
		return []byte{}, colsErr
	}
	buf.WriteString(string(colsJson))

	buf.WriteString(",\"rows\":")
	rowsJson, rowsErr := json.Marshal(grid.rows)
	if rowsErr != nil {
		return []byte{}, rowsErr
	}
	buf.WriteString(string(rowsJson))

	buf.WriteString("}")

	return []byte(buf.String()), nil

	// We can just marshal the struct here because of the struct type handling of json.Marshal along with the fact
	// that all fields of this struct are haystack standard and all the nested objects have a MarshalJSON method.
	// return json.Marshal(jsonGrid)
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

// WriteZincTo appends the Writer with the Grid zinc representation:
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

// Col is a column in a Grid.
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

// MarshalJSON represents the object in JSON object format with a required name field and the column metadata:
// "{"name":"<name>", "<field1>":<val1> ...}"
func (col *Col) MarshalJSON() ([]byte, error) {
	newMeta := col.meta.Set("name", NewStr(col.name))
	return json.Marshal(newMeta)
}

// WriteZincTo appends the Writer with the Col representation: <name> <meta>
func (col *Col) WriteZincTo(buf *bufio.Writer) {
	buf.WriteString(col.name)
	if !col.meta.IsEmpty() {
		buf.WriteRune(' ')
		col.meta.WriteZincTo(buf, false)
	}
}

// Row is a row in a Grid.
type Row struct {
	items map[string]Val
}

// Get returns the Val of the given name. If the name is not found, Null is returned.
func (row *Row) Get(name string) Val {
	return row.items[name]
}

// ToDict returns the values in a Dict format
func (row *Row) ToDict() Dict {
	return Dict{items: row.items}
}

// MarshalJSON represents the object in JSON object format: "{"<name1>":<val1>, "<name2>":<val2> ...}"
func (row *Row) MarshalJSON() ([]byte, error) {
	// Use Dict.MarshalJSON to enforce alphabetical order for easier testing.
	return json.Marshal(row.ToDict())
}

// WriteZincTo appends the Writer with the Row representation: <val1>, <val2>, ... Cols sets ordering
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
