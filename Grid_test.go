package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrid_ColCount(t *testing.T) {
	grid := newGridSimple()
	if grid.ColCount() != 2 {
		t.Error("List.Get returned an incorrect value")
	}
}

func TestGrid_Col(t *testing.T) {
	grid := newGridSimple()
	col := grid.Col("siteName")
	if col.Name() != "siteName" {
		t.Error("List.Col returned a column with the wrong name")
	}
	if col.Meta().ToZinc() != "{dis:\"Sites\"}" {
		t.Error("List.Col returned a column with the wrong meta")
	}
}

func TestGrid_ColAt(t *testing.T) {
	grid := newGridSimple()
	if grid.ColAt(0).Name() != "siteName" {
		t.Error("List.ColAt returned an incorrect value")
	}
}

func TestGrid_RenameCol(t *testing.T) {
	grid := newGridSimple()
	renamedGrid := grid.RenameCol("siteName", "siteName2")

	assert.Equal(t, grid.ColCount(), renamedGrid.ColCount())
	assert.Equal(t, grid.RowCount(), renamedGrid.RowCount())
	assert.NotNil(t, grid.Col("siteName"))
	assert.Nil(t, grid.Col("siteName2"))
	assert.Nil(t, renamedGrid.Col("siteName"))
	assert.NotNil(t, renamedGrid.Col("siteName2"))

	for i, row := range grid.Rows() {
		renamedRow := renamedGrid.RowAt(i)
		assert.Equal(t, row.Get("siteName"), renamedRow.Get("siteName2"))
	}
}

func TestGrid_RowCount(t *testing.T) {
	grid := newGridSimple()
	if grid.RowCount() != 2 {
		t.Error("List.RowCount returned an incorrect value")
	}
}

func TestGrid_RowAt(t *testing.T) {
	grid := newGridSimple()
	row := grid.RowAt(0)
	if row.ToDict().ToZinc() != "{siteName:\"Site 1\" val:356.214kW}" {
		t.Error("List.RowAt returned a row with the data")
	}
}

func TestGrid_ToZinc(t *testing.T) {
	grid := newGridSimple()
	zinc := `ver:"3.0" dis:"Site Energy Summary"
siteName dis:"Sites", val dis:"Value"
"Site 1", 356.214kW
"Site 2", 463.028kW`
	assert.Equal(t, grid.ToZinc(), zinc)
}

func TestGrid_ToZinc_empty(t *testing.T) {
	grid := EmptyGrid()
	zinc := `ver:"3.0"
empty
`
	assert.Equal(t, grid.ToZinc(), zinc)
}

func TestGrid_ToZinc_nested(t *testing.T) {
	grid := newGridNested()
	zinc := `ver:"3.0"
type, val
"list", [1, 2, 3]
"dict", {dis:"Dict!" foo}
"grid", <<
  ver:"3.0"
  a, b
  1, <<
    ver:"3.0"
    c, d
    5, 6
    >>
  3, 4
  >>
"scalar", "simple string"`
	assert.Equal(t, grid.ToZinc(), zinc)
}

func TestGrid_MarshalJSON(t *testing.T) {
	grid := newGridSimple()
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

func TestGrid_UnmarshalJSON(t *testing.T) {
	// Remember all dicts are alphabetical.
	jsonLiteral := "{" +
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
	zinc := `ver:"3.0" dis:"Site Energy Summary"
siteName dis:"Sites", val dis:"Value"
"Site 1", 356.214kW
"Site 2", 463.028kW`

	var grid Grid
	json.Unmarshal([]byte(jsonLiteral), &grid)
	assert.Equal(t, grid.ToZinc(), zinc)
}

func TestGrid_MarshalJSON_empty(t *testing.T) {
	grid := EmptyGrid()
	json := "{" +
		"\"meta\":{\"ver\":\"3.0\"}," +
		"\"cols\":[]," + // Empty lists are marshaled as null
		"\"rows\":[]" +
		"}"
	valTest_MarshalJSON(grid, json, t)
}

func TestGrid_MarshalJSON_nested(t *testing.T) {
	grid := newGridNested()
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

func TestGrid_MarshalHayson(t *testing.T) {
	grid := newGridSimple()
	// Remember all dicts are alphabetical.
	json := "{" +
		"\"_kind\":\"grid\"," +
		"\"meta\":{\"_kind\":\"dict\",\"dis\":\"Site Energy Summary\",\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"_kind\":\"dict\",\"dis\":\"Sites\",\"name\":\"siteName\"}," +
		"{\"_kind\":\"dict\",\"dis\":\"Value\",\"name\":\"val\"}" +
		"]," +
		"\"rows\":[" +
		"{\"_kind\":\"dict\",\"siteName\":\"Site 1\",\"val\":{\"_kind\":\"number\",\"val\":356.214,\"unit\":\"kW\"}}," +
		"{\"_kind\":\"dict\",\"siteName\":\"Site 2\",\"val\":{\"_kind\":\"number\",\"val\":463.028,\"unit\":\"kW\"}}" +
		"]" +
		"}"
	valTest_MarshalHayson(grid, json, t)
}

func TestGrid_MarshalHayson_empty(t *testing.T) {
	grid := EmptyGrid()
	json := "{" +
		"\"_kind\":\"grid\"," +
		"\"meta\":{\"_kind\":\"dict\",\"ver\":\"3.0\"}," +
		"\"cols\":[]," +
		"\"rows\":[]" +
		"}"
	valTest_MarshalHayson(grid, json, t)
}

func TestGrid_MarshalHayson_nested(t *testing.T) {
	grid := newGridNested()
	// I tried to make it easier to see, but the go formatting is making it hard. See raw value below and use a formatter if needed
	// {"_kind":"grid","meta":{"_kind":"dict","ver":"3.0"},"cols":[{"_kind":"dict","name":"type"},{"_kind":"dict","name":"val"}],"rows":[{"_kind":"dict","type":"list","val":[{"_kind":"number","val":1},{"_kind":"number","val":2},{"_kind":"number","val":3}]},{"_kind":"dict","type":"dict","val":{"_kind":"dict","dis":"Dict!","foo":{"_kind":"marker"}}},{"_kind":"dict","type":"grid","val":{"_kind":"grid","meta":{"_kind":"dict","ver":"3.0"},"cols":[{"_kind":"dict","name":"a"},{"_kind":"dict","name":"b"}],"rows":[{"_kind":"dict","a":{"_kind":"number","val":1},"b":{"_kind":"grid","meta":{"_kind":"dict","ver":"3.0"},"cols":[{"_kind":"dict","name":"c"},{"_kind":"dict","name":"d"}],"rows":[{"_kind":"dict","c":{"_kind":"number","val":5},"d":{"_kind":"number","val":6}}]}},{"_kind":"dict","a":{"_kind":"number","val":3},"b":{"_kind":"number","val":4}}]}},{"_kind":"dict","type":"scalar","val":"simple string"}]}
	json := "{" +
		"\"_kind\":\"grid\"," +
		"\"meta\":{\"_kind\":\"dict\",\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"_kind\":\"dict\",\"name\":\"type\"}," +
		"{\"_kind\":\"dict\",\"name\":\"val\"}" +
		"]," +
		"\"rows\":[" +
		"{\"_kind\":\"dict\",\"type\":\"list\",\"val\":[{\"_kind\":\"number\",\"val\":1},{\"_kind\":\"number\",\"val\":2},{\"_kind\":\"number\",\"val\":3}]}," +
		"{\"_kind\":\"dict\",\"type\":\"dict\",\"val\":{\"_kind\":\"dict\",\"dis\":\"Dict!\",\"foo\":{\"_kind\":\"marker\"}}}," +
		"{\"_kind\":\"dict\",\"type\":\"grid\",\"val\":{" + // Start nested 1
		"\"_kind\":\"grid\"," +
		"\"meta\":{\"_kind\":\"dict\",\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"_kind\":\"dict\",\"name\":\"a\"}," +
		"{\"_kind\":\"dict\",\"name\":\"b\"}" +
		"]," +
		"\"rows\":[" +
		"{\"_kind\":\"dict\",\"a\":{\"_kind\":\"number\",\"val\":1},\"b\":{" + // Start nested 2
		"\"_kind\":\"grid\"," +
		"\"meta\":{\"_kind\":\"dict\",\"ver\":\"3.0\"}," +
		"\"cols\":[" +
		"{\"_kind\":\"dict\",\"name\":\"c\"}," +
		"{\"_kind\":\"dict\",\"name\":\"d\"}" +
		"]," +
		"\"rows\":[" +
		"{\"_kind\":\"dict\",\"c\":{\"_kind\":\"number\",\"val\":5},\"d\":{\"_kind\":\"number\",\"val\":6}}" +
		"]" +
		"}" + // End nested 2
		"}," +
		"{\"_kind\":\"dict\",\"a\":{\"_kind\":\"number\",\"val\":3},\"b\":{\"_kind\":\"number\",\"val\":4}}" +
		"]" +
		"}" + // End nested 1
		"}," +
		"{\"_kind\":\"dict\",\"type\":\"scalar\",\"val\":\"simple string\"}" +
		"]" +
		"}"
	valTest_MarshalHayson(grid, json, t)
}

// Zinc representation:
//
//	ver:"3.0" dis:"Site Energy Summary"
//	siteName dis:"Sites", val dis:"Value"
//	"Site 1", 356.214kW
//	"Site 2", 463.028kW
func newGridSimple() Grid {
	gb := NewGridBuilder()
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
	return gb.ToGrid()
}

// Zinc representation:
// ver:"3.0"
// type, val
// "list", [1, 2, 3]
// "dict", {dis:"Dict!" foo}
// "grid", <<
//
//	ver:"3.0"
//	a, b
//	1, <<
//	  ver:"3.0"
//	  c, d
//	  5, 6
//	  >>
//	3, 4
//	>>
//
// "scalar", "simple string"
func newGridNested() Grid {
	gb := NewGridBuilder()
	gb.AddColNoMeta("type")
	gb.AddColNoMeta("val")
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
	dblNestedGb.AddColNoMeta("c")
	dblNestedGb.AddColNoMeta("d")
	dblNestedGb.AddRow(
		[]Val{
			NewNumber(5, ""),
			NewNumber(6, ""),
		},
	)
	dblNestedGrid := dblNestedGb.ToGrid()
	var nestedGb GridBuilder
	nestedGb.AddColNoMeta("a")
	nestedGb.AddColNoMeta("b")
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
	return gb.ToGrid()
}
