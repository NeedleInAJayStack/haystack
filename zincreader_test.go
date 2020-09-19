package haystack

import (
	"testing"
)

func TestZincReader_empty(t *testing.T) {
	input := "ver:\"3.0\" tag:N\n" +
		"a nullmetatag:N, b markermetatag\n" +
		""

	var gb GridBuilder
	gb.SetMeta(
		map[string]Val{
			"tag": NewNull(),
		},
	)
	gb.AddColWMeta(
		"a",
		map[string]Val{
			"nullmetatag": NewNull(),
		},
	)
	gb.AddColWMeta(
		"b",
		map[string]Val{
			"markermetatag": NewMarker(),
		},
	)
	expected := gb.ToGrid()
	testZincReaderGrid(t, input, expected)
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
		t.Error("Grids do not match\nACTUAL:\n" + actualZinc + "\nEXPECTED:\n" + expectedZinc)
	}
}
