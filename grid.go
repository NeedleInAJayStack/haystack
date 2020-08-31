package haystack

type Col struct {
	index int
	name  string
	meta  Dict
}

type Row struct {
	vals []Val
}

type Grid struct {
	meta Dict
	cols []Col
	rows []Row
}

// Format as <name> <meta>
func (col *Col) encode() string {
	return col.name + " " + col.meta.encode(false)
}

// Format as <val1>, <val2>, ...
func (row *Row) encode() string {
	result := ""
	for idx, val := range row.vals {
		if idx != 0 {
			result = result + ", "
		}
		result = result + val.toZinc()
	}
	return result
}

// Format as:
//     ver:"3.0" <meta>
//     <col1>, <col2>, ...
//     <row1>
//     <row2>
//     ...
func (grid *Grid) toZinc() string {
	result := "ver:\"3.0\""
	if !grid.meta.isEmpty() {
		result = result + " " + grid.meta.encode(false)
	}
	result = result + "\n"
	for colIdx, col := range grid.cols {
		if colIdx != 0 {
			result = result + ", "
		}
		result = result + col.encode()
	}
	result = result + "\n"
	for rowIdx, row := range grid.rows {
		if rowIdx != 0 {
			result = result + "\n"
		}
		result = result + row.encode()
	}
	return result
}
