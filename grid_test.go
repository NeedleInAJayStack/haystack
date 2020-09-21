package haystack

import "testing"

func TestGrid_ToZinc(t *testing.T) {
	var gb GridBuilder
	gb.SetMeta(
		map[string]Val{
			"dis": NewStr("Site Energy Summary"),
		},
	)
	gb.AddCol(
		"siteName",
		map[string]Val{
			"dis": NewStr("Sites"),
		},
	)
	gb.AddCol(
		"val",
		map[string]Val{
			"dis": NewStr("Value"),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("Site 1"),
			NewNumber(356.214, "kW"),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("Site 2"),
			NewNumber(463.028, "kW"),
		},
	)
	actual := gb.ToGrid().ToZinc()
	expected := "ver:\"3.0\" dis:\"Site Energy Summary\"\n" +
		"siteName dis:\"Sites\", val dis:\"Value\"\n" +
		"\"Site 1\", 356.214kW\n" +
		"\"Site 2\", 463.028kW"
	if actual != expected {
		t.Error("Grids don't match\nACTUAL:\n" + actual + "\nEXPECTED:\n" + expected)
	}
}

func TestGrid_ToZinc_nested(t *testing.T) {
	var gb GridBuilder
	gb.AddCol("type", map[string]Val{})
	gb.AddCol("val", map[string]Val{})
	gb.AddRow(
		[]Val{
			NewStr("list"),
			NewList(
				[]Val{
					NewNumber(1, ""),
					NewNumber(2, ""),
					NewNumber(3, ""),
				},
			),
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("dict"),
			NewDict(
				map[string]Val{
					"dis": NewStr("Dict!"),
					"foo": NewMarker(),
				},
			),
		},
	)
	var dblNestedGb GridBuilder
	dblNestedGb.AddCol("c", map[string]Val{})
	dblNestedGb.AddCol("d", map[string]Val{})
	dblNestedGb.AddRow(
		[]Val{
			NewNumber(5, ""),
			NewNumber(6, ""),
		},
	)
	dblNestedGrid := dblNestedGb.ToGrid()
	var nestedGb GridBuilder
	nestedGb.AddCol("a", map[string]Val{})
	nestedGb.AddCol("b", map[string]Val{})
	nestedGb.AddRow(
		[]Val{
			NewNumber(1, ""),
			dblNestedGrid,
		},
	)
	nestedGb.AddRow(
		[]Val{
			NewNumber(3, ""),
			NewNumber(4, ""),
		},
	)
	nestedGrid := nestedGb.ToGrid()
	gb.AddRow(
		[]Val{
			NewStr("grid"),
			nestedGrid,
		},
	)
	gb.AddRow(
		[]Val{
			NewStr("scalar"),
			NewStr("simple string"),
		},
	)
	actual := gb.ToGrid().ToZinc()
	expected := "ver:\"3.0\"\n" +
		"type, val\n" +
		"\"list\", [1, 2, 3]\n" +
		"\"dict\", {dis:\"Dict!\" foo}\n" +
		"\"grid\", <<\n" +
		"  ver:\"3.0\"\n" +
		"  a, b\n" +
		"  1, <<\n" +
		"    ver:\"3.0\"\n" +
		"    c, d\n" +
		"    5, 6\n" +
		"    >>\n" +
		"  3, 4\n" +
		"  >>\n" +
		"\"scalar\", \"simple string\""
	if actual != expected {
		t.Error("Grids don't match\nACTUAL:\n" + actual + "\nEXPECTED:\n" + expected)
	}
}
