package haystack

import "testing"

func TestGrid_ToZinc(t *testing.T) {
	grid := Grid{
		meta: Dict{
			items: map[string]Val{
				"dis": &Str{val: "Site Energy Summary"},
			},
		},
		cols: []Col{
			Col{
				name: "siteName",
				meta: Dict{
					items: map[string]Val{
						"dis": &Str{val: "Sites"},
					},
				},
			},
			Col{
				name: "val",
				meta: Dict{
					items: map[string]Val{
						"dis": &Str{val: "Value"},
					},
				},
			},
		},
		rows: []Row{
			Row{
				vals: []Val{
					&Str{val: "Site 1"},
					&Number{val: 356.214, unit: "kW"},
				},
			},
			Row{
				vals: []Val{
					&Str{val: "Site 2"},
					&Number{val: 463.028, unit: "kW"},
				},
			},
		},
	}
	gridZinc := grid.ToZinc()
	expected := "ver:\"3.0\" dis:\"Site Energy Summary\"\n" +
		"siteName dis:\"Sites\", val dis:\"Value\"\n" +
		"\"Site 1\", 356.214kW\n" +
		"\"Site 2\", 463.028kW"
	if gridZinc != expected {
		t.Error(gridZinc)
	}
}
