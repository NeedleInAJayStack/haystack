package haystack

import "testing"

func TestGrid_ToZinc(t *testing.T) {
	var gb GridBuilder
	gb.SetMeta(
		map[string]Val{
			"dis": NewStr("Site Energy Summary"),
		},
	)
	gb.AddColWMeta(
		"siteName",
		map[string]Val{
			"dis": NewStr("Sites"),
		},
	)
	gb.AddColWMeta(
		"val",
		map[string]Val{
			"dis": NewStr("Value"),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("Site 1"),
			&Number{val: 356.214, unit: "kW"},
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("Site 2"),
			&Number{val: 463.028, unit: "kW"},
		},
	)
	gridZinc := gb.ToGrid().ToZinc()
	expected := "ver:\"3.0\" dis:\"Site Energy Summary\"\n" +
		"siteName dis:\"Sites\", val dis:\"Value\"\n" +
		"\"Site 1\", 356.214kW\n" +
		"\"Site 2\", 463.028kW"
	if gridZinc != expected {
		t.Error(gridZinc)
	}
}
