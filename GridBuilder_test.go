package haystack

import (
	"testing"
)

func TestGridBuilder_ToGrid(t *testing.T) {
	gb := NewGridBuilder()
	gb.AddColNoMeta("col1")
	gb.AddRow([]Val{
		NewStr("val1"),
	})
	grid1 := gb.ToGrid()

	gb.AddRow([]Val{
		NewStr("val2"),
	})
	grid2 := gb.ToGrid()

	if grid1.ToZinc() == grid2.ToZinc() {
		t.Error("GridBuilder.ToGrid persists grid values")
	}
}
