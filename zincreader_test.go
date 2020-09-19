package haystack

import (
	"testing"
)

func TestZincReader_empty(t *testing.T) {
	testZincReaderGrid(
		t,
		"ver:\"3.0\" tag:N\n"+
			"a nullmetatag:N, b markermetatag\n"+
			"",
		Grid{
			meta: Dict{
				items: map[string]Val{
					"tag": &Null{},
				},
			},
			cols: []Col{
				Col{
					name: "a",
					meta: Dict{
						items: map[string]Val{
							"nullmetatag": NewNull(),
						},
					},
				},
				Col{
					name: "b",
					meta: Dict{
						items: map[string]Val{
							"markermetatag": NewMarker(),
						},
					},
				},
			},
			rows: []Row{},
		},
	)
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'ToZinc' method
func testZincReaderGrid(t *testing.T, str string, expected Grid) {
	var reader ZincReader
	reader.InitString(str)

	val := reader.ReadVal()
	actual := val.(Grid)

	testGridEq(t, actual, expected)
}

// Test whether the grids match based on a ToZinc call
func testGridEq(t *testing.T, actual Grid, expected Grid) {
	actualZinc := actual.ToZinc()
	expectedZinc := expected.ToZinc()

	if actualZinc != expectedZinc {
		t.Error("Grids do not match\n" + "ACTUAL:\n" + actualZinc + "EXPECTED:\n" + expectedZinc)
	}
}
