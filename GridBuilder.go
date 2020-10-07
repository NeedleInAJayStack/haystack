package haystack

// GridBuilder is used to easily construct a Grid instance.
type GridBuilder struct {
	meta Dict
	cols []Col
	rows []Row
}

// ToGrid returns the grid representation of the builder.
func (gb *GridBuilder) ToGrid() Grid {
	return Grid{
		meta: gb.meta,
		cols: gb.cols,
		rows: gb.rows,
	}
}

// SetMeta erases any existing meta and replaces it with the input map.
func (gb *GridBuilder) SetMeta(meta map[string]Val) {
	gb.SetMetaDict(NewDict(meta))
}

// SetMetaDict erases any existing meta and replaces it with the input Dict.
func (gb *GridBuilder) SetMetaDict(meta Dict) {
	gb.meta = meta
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
