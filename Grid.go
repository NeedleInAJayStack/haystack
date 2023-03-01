package haystack

import (
	"bufio"
	"encoding/json"
	"errors"
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
	return Grid{
		meta: Dict{},
		cols: []Col{},
		rows: []Row{},
	}
}

// Meta returns the grid-level metadata
func (grid Grid) Meta() Dict {
	return grid.meta
}

// ColCount returns the count of columns
func (grid Grid) ColCount() int {
	return len(grid.cols)
}

// Cols returns the column objects
func (grid Grid) Cols() []Col {
	return grid.cols
}

// Col returns the column matching the name
func (grid Grid) Col(name string) Col {
	var colMatch Col
	for i := 0; i < grid.ColCount(); i++ {
		col := grid.ColAt(i)
		if grid.cols[i].name == name {
			colMatch = col
			return colMatch
		}
	}
	panic("Unknown column: " + name)
}

// ColAt returns the column at the index
func (grid Grid) ColAt(index int) Col {
	return grid.cols[index]
}

// RowCount returns the count of rows
func (grid Grid) RowCount() int {
	return len(grid.rows)
}

// Rows returns the column objects
func (grid Grid) Rows() []Row {
	return grid.rows
}

// RowAt returns the row at the index
func (grid Grid) RowAt(index int) Row {
	return grid.rows[index]
}

// MarshalJSON represents the object in a special JSON object format. See https://project-haystack.org/doc/Json#grid
func (grid Grid) MarshalJSON() ([]byte, error) {
	buf := strings.Builder{}

	buf.WriteString("{\"meta\":")
	newMeta := grid.meta.Set("ver", NewStr("3.0")) // Add in version
	metaJSON, metaErr := json.Marshal(newMeta)
	if metaErr != nil {
		return []byte{}, metaErr
	}
	buf.Write(metaJSON)

	buf.WriteString(",\"cols\":")
	colsJSON, colsErr := json.Marshal(grid.cols)
	if colsErr != nil {
		return []byte{}, colsErr
	}
	buf.Write(colsJSON)

	buf.WriteString(",\"rows\":")
	rowsJSON, rowsErr := json.Marshal(grid.rows)
	if rowsErr != nil {
		return []byte{}, rowsErr
	}
	buf.Write(rowsJSON)

	buf.WriteString("}")

	return []byte(buf.String()), nil

	// We can just marshal the struct here because of the struct type handling of json.Marshal along with the fact
	// that all fields of this struct are haystack standard and all the nested objects have a MarshalJSON method.
	// return json.Marshal(jsonGrid)
}

// UnmarshalJSON interprets the special JSON object format. See https://project-haystack.org/doc/Json#grid
func (grid *Grid) UnmarshalJSON(buf []byte) error {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(buf, &jsonMap)
	if err != nil {
		return err
	}

	newGrid, newErr := gridFromJSON(jsonMap)
	*grid = newGrid
	return newErr
}

func gridFromJSON(jsonMap map[string]interface{}) (Grid, error) {
	gb := NewGridBuilder()

	if jsonMap["meta"] == nil || jsonMap["cols"] == nil || jsonMap["rows"] == nil {
		return EmptyGrid(), errors.New("object does not contain keys: meta, cols, or rows")
	}

	metaMap := jsonMap["meta"].(map[string]interface{})
	delete(metaMap, "ver")
	meta, metaErr := dictFromJSON(metaMap)
	if metaErr != nil {
		return EmptyGrid(), metaErr
	}
	gb.AddMetaDict(meta)

	for _, jsonCol := range jsonMap["cols"].([]interface{}) {
		jsonColMap := jsonCol.(map[string]interface{})

		name := jsonColMap["name"].(string)
		delete(jsonColMap, "name")
		colMeta, colMetaErr := dictFromJSON(jsonColMap)
		if colMetaErr != nil {
			return EmptyGrid(), colMetaErr
		}
		gb.AddColDict(name, colMeta)
	}

	for _, jsonRow := range jsonMap["rows"].([]interface{}) {
		jsonRowMap := jsonRow.(map[string]interface{})
		row, rowErr := dictFromJSON(jsonRowMap)
		if rowErr != nil {
			return EmptyGrid(), rowErr
		}
		gb.AddRowDict(row)
	}

	return gb.ToGrid(), nil
}

// MarshalHayson represents the object in a special JSON object format. See https://bitbucket.org/finproducts/hayson/src/master/spec.md
func (grid Grid) MarshalHayson() ([]byte, error) {
	buf := strings.Builder{}

	buf.WriteString("{\"_kind\":\"grid\",\"meta\":")
	newMeta := grid.meta.Set("ver", NewStr("3.0")) // Add in version
	metaBytes, metaErr := newMeta.MarshalHayson()
	if metaErr != nil {
		return []byte{}, metaErr
	}
	buf.Write(metaBytes)

	buf.WriteString(",\"cols\":[")
	for idx, col := range grid.cols {
		if idx != 0 {
			buf.WriteString(",")
		}
		colHayson, colErr := col.MarshalHayson()
		if colErr != nil {
			return []byte{}, colErr
		}
		buf.Write(colHayson)
	}
	buf.WriteString("]")

	buf.WriteString(",\"rows\":[")
	for idx, row := range grid.rows {
		if idx != 0 {
			buf.WriteString(",")
		}
		rowHayson, rowErr := row.MarshalHayson()
		if rowErr != nil {
			return []byte{}, rowErr
		}
		buf.Write(rowHayson)
	}
	buf.WriteString("]")
	buf.WriteString("}")

	return []byte(buf.String()), nil
}

// ToZinc representes the object as:
//
//	ver:"3.0" <meta>
//	<col1>, <col2>, ...
//	<row1>
//	<row2>
//	...
func (grid Grid) ToZinc() string {
	builder := new(strings.Builder)
	out := bufio.NewWriter(builder)
	grid.WriteZincTo(out, 0)
	out.Flush()
	return builder.String()
}

// WriteZincTo appends the Writer with the Grid zinc representation:
//
//	ver:"3.0" <meta>
//	<col1>, <col2>, ...
//	<row1>
//	<row2>
//	...
//
// indentSize is the number of spaces to add to each new-line.
func (grid Grid) WriteZincTo(buf *bufio.Writer, indentSize int) {
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
			row.WriteZincTo(buf, grid.cols, indentSize)
		}
	}
}

// Col is a column in a Grid.
type Col struct {
	index int
	name  string
	meta  Dict
}

func newCol(index int, name string, meta Dict) Col {
	return Col{
		index: index,
		name:  name,
		meta:  meta,
	}
}

// Name returns the string name
func (col Col) Name() string {
	return col.name
}

// Meta returns the column-level metadata
func (col Col) Meta() Dict {
	return col.meta
}

// MarshalJSON represents the object in JSON object format with a required name field and the column metadata:
// "{"name":"<name>", "<field1>":<val1> ...}"
func (col Col) MarshalJSON() ([]byte, error) {
	newMeta := col.meta.Set("name", NewStr(col.name))
	return newMeta.MarshalJSON()
}

// MarshalHayson represents the object in JSON object format with a required name field and the column metadata:
// "{"name":"<name>", "<field1>":<val1> ...}"
func (col Col) MarshalHayson() ([]byte, error) {
	newMeta := col.meta.Set("name", NewStr(col.name))
	return newMeta.MarshalHayson()
}

// WriteZincTo appends the Writer with the Col representation: <name> <meta>
func (col Col) WriteZincTo(buf *bufio.Writer) {
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
func (row Row) Get(name string) Val {
	return row.items[name]
}

// ToDict returns the values in a Dict format
func (row Row) ToDict() Dict {
	return NewDict(row.items)
}

// MarshalJSON represents the object in JSON object format: "{"<name1>":<val1>, "<name2>":<val2> ...}"
func (row Row) MarshalJSON() ([]byte, error) {
	// Use Dict.MarshalJSON to enforce alphabetical order for easier testing.
	return row.ToDict().MarshalJSON()
}

// MarshalHayson represents the object in JSON object format: "{"<name1>":<val1>, "<name2>":<val2> ...}"
func (row Row) MarshalHayson() ([]byte, error) {
	// Use Dict.MarshalHayson to enforce alphabetical order for easier testing.
	return row.ToDict().MarshalHayson()
}

// WriteZincTo appends the Writer with the Row representation: <val1>, <val2>, ... Cols sets ordering
func (row Row) WriteZincTo(buf *bufio.Writer, cols []Col, indentSize int) {
	for colIdx, col := range cols {
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
