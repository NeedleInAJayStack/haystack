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
							"nullmetatag": &Null{},
						},
					},
				},
				Col{
					name: "b",
					meta: Dict{
						items: map[string]Val{
							"markermetatag": &Marker{},
						},
					},
				},
			},
			rows: []Row{},
		},
	)
}

// Verifies that the tokenized result has the expected token type and value.
// Values are matched based on the result of the 'toZinc' method
func testZincReaderGrid(t *testing.T, str string, expected Grid) {
	var reader ZincReader
	reader.InitString(str)

	val, err := reader.ReadVal()
	actual := val.(Grid)
	if err != nil {
		t.Error(err)
	}

	testGridEq(t, actual, expected)
}

// Test whether the grids match based on a toZinc call
func testGridEq(t *testing.T, actual Grid, expected Grid) {
	actualZinc := actual.toZinc()
	expectedZinc := expected.toZinc()

	if actualZinc != expectedZinc {
		t.Error("Grids do not match")
	}
}
