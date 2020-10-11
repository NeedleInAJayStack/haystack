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
	grid := gb.ToGrid()
	zinc := "ver:\"3.0\" dis:\"Site Energy Summary\"\n" +
		"siteName dis:\"Sites\", val dis:\"Value\"\n" +
		"\"Site 1\", 356.214kW\n" +
		"\"Site 2\", 463.028kW"
	valTest_ToZinc_Grid(grid, zinc, t)
}

func TestGrid_ToZinc_empty(t *testing.T) {
	grid := Grid{}
	zinc := "ver:\"3.0\"\n" +
		"empty\n"
	valTest_ToZinc_Grid(grid, zinc, t)
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
	grid := gb.ToGrid()
	zinc := "ver:\"3.0\"\n" +
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
	valTest_ToZinc_Grid(grid, zinc, t)
}

func TestGrid_MarshalJSON(t *testing.T) {
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
	grid := gb.ToGrid()
	// Remember all dicts are alphabetical.
	json := "{" +
		"\"meta\":{\"dis\":\"Site Energy Summary\",\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"dis\":\"Sites\",\"name\":\"siteName\"}," +
		"{\"dis\":\"Value\",\"name\":\"val\"}" +
		"]," +
		"\"rows\":[" +
		"{\"siteName\":\"Site 1\",\"val\":\"n:356.214 kW\"}," +
		"{\"siteName\":\"Site 2\",\"val\":\"n:463.028 kW\"}" +
		"]" +
		"}"
	valTest_MarshalJSON(grid, json, t)
}

func TestGrid_MarshalJSON_empty(t *testing.T) {
	grid := Grid{}
	json := "{" +
		"\"meta\":{\"ver\":\"3.0\"}," +
		"\"cols\":null," + // Empty lists are marshaled as null
		"\"rows\":null" +
		"}"
	valTest_MarshalJSON(grid, json, t)
}

func TestGrid_MarshalJSON_nested(t *testing.T) {
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
	grid := gb.ToGrid()
	// I tried to make it easier to see, but the go formatting is making it hard. See raw value below and use a formatter if needed
	// {"meta":{"ver":"3.0"},"cols":[{"name":"type"},{"name":"val"}],"rows":[{"type":"list","val":["n:1","n:2","n:3"]},{"type":"dict","val":{"dis":"Dict!","foo":"m:"}},{"type":"grid","val":{"meta":{"ver":"3.0"},"cols":[{"name":"a"},{"name":"b"}],"rows":[{"a":"n:1","b":{"meta":{"ver":"3.0"},"cols":[{"name":"c"},{"name":"d"}],"rows":[{"c":"n:5","d":"n:6"}]}},{"a":"n:3","b":"n:4"}]}},{"type":"scalar","val":"simple string"}]}
	json := "{" +
		"\"meta\":{\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"name\":\"type\"}," +
		"{\"name\":\"val\"}" +
		"]," +
		"\"rows\":[" +
		"{\"type\":\"list\",\"val\":[\"n:1\",\"n:2\",\"n:3\"]}," +
		"{\"type\":\"dict\",\"val\":{\"dis\":\"Dict!\",\"foo\":\"m:\"}}," +
		"{\"type\":\"grid\",\"val\":{" + // Start nested 1
		"\"meta\":{\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"name\":\"a\"}," +
		"{\"name\":\"b\"}" +
		"]," +
		"\"rows\":[" +
		"{\"a\":\"n:1\",\"b\":{" + // Start nested 2
		"\"meta\":{\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"name\":\"c\"}," +
		"{\"name\":\"d\"}" +
		"]," +
		"\"rows\":[" +
		"{\"c\":\"n:5\",\"d\":\"n:6\"}" +
		"]" +
		"}" + // End nested 2
		"}," +
		"{\"a\":\"n:3\",\"b\":\"n:4\"}" +
		"]" +
		"}" + // End nested 1
		"}," +
		"{\"type\":\"scalar\",\"val\":\"simple string\"}" +
		"]" +
		"}"
	valTest_MarshalJSON(grid, json, t)
}
