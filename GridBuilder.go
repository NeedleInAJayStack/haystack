package haystack

// GridBuilder is used to easily construct a Grid instance.
type GridBuilder struct {
	meta map[string]Val
	cols []Col
	rows []Row
}

func NewGridBuilder() *GridBuilder {
	return &GridBuilder{
		meta: map[string]Val{},
		cols: []Col{},
		rows: []Row{},
	}
}

// ToGrid returns the grid representation of the builder.
func (gb *GridBuilder) ToGrid() *Grid {
	meta := NewDict(gb.meta)
	return &Grid{
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

// AddMetaDict adds or replaces the meta keys with the inputs.
func (gb *GridBuilder) AddMetaDict(meta *Dict) {
	gb.AddMeta(meta.items)
}

// AddMeta adds or replaces the meta keys with the inputs.
func (gb *GridBuilder) AddMetaVal(name string, val Val) {
	gb.meta[name] = val
}

// SetMeta erases any existing meta and replaces it with the input mappings.
func (gb *GridBuilder) SetMeta(meta map[string]Val) {
	gb.meta = map[string]Val{} // Empty meta values
	gb.AddMeta(meta)
}

// SetMetaDict erases any existing meta and replaces it with the input Dict.
func (gb *GridBuilder) SetMetaDict(meta *Dict) {
	gb.SetMeta(meta.items)
}

// AddCol adds a column with the given name and meta map.
func (gb *GridBuilder) AddCol(name string, meta map[string]Val) {
	gb.AddColDict(name, NewDict(meta))
}

// AddColNoMeta adds a column with the given name and empty meta.
func (gb *GridBuilder) AddColNoMeta(name string) {
	gb.AddColDict(name, EmptyDict())
}

// AddColDict adds a column with the given name and meta Dict.
func (gb *GridBuilder) AddColDict(name string, meta *Dict) {
	index := len(gb.cols)
	newCol := NewCol(index, name, meta)
	// TODO check that the name doesn't duplicate
	gb.cols = append(gb.cols, newCol)
}

// AddColMeta adds the metadata to an existing column with the given name.
func (gb *GridBuilder) AddColMeta(name string, meta map[string]Val) {
	col := gb.getCol(name)
	col.meta = col.meta.SetAll(meta)
}

// AddColMetaVal adds the metadata name and value to an existing column with the given name.
func (gb *GridBuilder) AddColMetaVal(colName string, metaName string, metaVal Val) {
	gb.AddColMeta(colName, map[string]Val{metaName: metaVal})
}

// AddColMetaDict adds the metadata Dict to an existing column with the given name.
func (gb *GridBuilder) AddColMetaDict(name string, meta *Dict) {
	gb.AddColMeta(name, meta.items)
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
func (gb *GridBuilder) AddRowDict(row *Dict) {
	items := make(map[string]Val)
	for _, col := range gb.cols {
		items[col.name] = row.Get(col.name)
	}
	newRow := Row{items: items}
	gb.rows = append(gb.rows, newRow)
}

// AddRowDicts adds rows from the dicts, extracting the values that correspond to the grid columns.
func (gb *GridBuilder) AddRowDicts(rows []*Dict) {
	for _, row := range rows {
		gb.AddRowDict(row)
	}
}

func (gb *GridBuilder) getCol(colName string) Col {
	var matchCol Col
	match := false
	for _, col := range gb.cols {
		if col.name == colName {
			matchCol = col
			match = true
		}
	}
	if !match {
		return matchCol
	}
	panic("Col with name not found: " + colName)
}
