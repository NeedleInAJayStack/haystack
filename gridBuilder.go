package haystack

type GridBuilder struct {
	meta Dict
	cols []Col
	rows []Row
}

func (gb *GridBuilder) ToGrid() Grid {
	return Grid{
		meta: gb.meta,
		cols: gb.cols,
		rows: gb.rows,
	}
}

func (gb *GridBuilder) SetMeta(meta map[string]Val) {
	gb.SetMetaDict(NewDict(meta))
}

func (gb *GridBuilder) SetMetaDict(meta Dict) {
	gb.meta = meta
}

func (gb *GridBuilder) AddCol(name string, meta map[string]Val) {
	gb.AddColDict(name, NewDict(meta))
}

func (gb *GridBuilder) AddColNoMeta(name string) {
	gb.AddColDict(name, NewDict(map[string]Val{}))
}

func (gb *GridBuilder) AddColDict(name string, meta Dict) {
	index := len(gb.cols)
	newCol := Col{index: index, name: name, meta: meta}
	// TODO check that the name doesn't duplicate
	gb.cols = append(gb.cols, newCol)
}

func (gb *GridBuilder) AddRow(vals []Val) {
	items := make(map[string]Val)
	for idx, col := range gb.cols {
		items[col.name] = vals[idx]
	}
	newRow := Row{items: items}
	gb.rows = append(gb.rows, newRow)
}
