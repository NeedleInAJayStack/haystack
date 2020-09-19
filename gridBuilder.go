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
	gb.meta = NewDict(meta)
}

func (gb *GridBuilder) AddCol(name string) {
	newCol := Col{name: name}
	gb.cols = append(gb.cols, newCol)
}

func (gb *GridBuilder) AddColWMeta(name string, meta map[string]Val) {
	newCol := Col{name: name, meta: NewDict(meta)}
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
