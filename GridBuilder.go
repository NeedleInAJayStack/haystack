package haystack

// GridBuilder is used to easily construct a Grid instance.
type GridBuilder struct {
	meta map[string]Val
	cols []Col
	rows []Row
}

// ToGrid returns the grid representation of the builder.
func (gb *GridBuilder) ToGrid() Grid {
	meta := NewDict(gb.meta)
	return Grid{
		meta: meta,
		cols: gb.cols,
		rows: gb.rows,
	}
}

// AddMeta adds or replaces the meta keys with the inputs.
func (gb *GridBuilder) AddMeta(meta map[string]Val) {
	for name, val := range meta {
		gb.AddMetaVal(name, val)
	}
}

// AddMeta adds or replaces the meta keys with the inputs.
func (gb *GridBuilder) AddMetaVal(name string, val Val) {
	gb.meta[name] = val
}

// SetMeta erases any existing meta and replaces it with the input map.
func (gb *GridBuilder) SetMeta(meta map[string]Val) {
	gb.meta = meta
}

// SetMetaDict erases any existing meta and replaces it with the input Dict.
func (gb *GridBuilder) SetMetaDict(meta Dict) {
	gb.meta = meta.items
}

// AddCol adds a column with the given name and meta map.
func (gb *GridBuilder) AddCol(name string, meta map[string]Val) {
	gb.AddColDict(name, NewDict(meta))
}

// AddColNoMeta adds a column with the given name and empty meta.
func (gb *GridBuilder) AddColNoMeta(name string) {
	gb.AddColDict(name, NewDict(map[string]Val{}))
}

// AddColDict adds a column with the given name and meta Dict.
func (gb *GridBuilder) AddColDict(name string, meta Dict) {
	index := len(gb.cols)
	newCol := Col{index: index, name: name, meta: meta}
	// TODO check that the name doesn't duplicate
	gb.cols = append(gb.cols, newCol)
}

// AddRow adds a row with the input values, according to the column order.
func (gb *GridBuilder) AddRow(vals []Val) {
	items := make(map[string]Val)
	for idx, col := range gb.cols {
		items[col.name] = vals[idx]
	}
	newRow := Row{items: items}
	gb.rows = append(gb.rows, newRow)
}

// AddRowDict adds a row from the input dict, extracting the values that correspond to the grid columns.
func (gb *GridBuilder) AddRowDict(row Dict) {
	items := make(map[string]Val)
	for _, col := range gb.cols {
		items[col.name] = row.Get(col.name)
	}
	newRow := Row{items: items}
	gb.rows = append(gb.rows, newRow)
}

// AddRowDicts adds rows from the dicts, extracting the values that correspond to the grid columns.
func (gb *GridBuilder) AddRowDicts(rows []Dict) {
	for _, row := range rows {
		gb.AddRowDict(row)
	}
}
