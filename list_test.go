package haystack

import "testing"

func TestList_ToZinc(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	listZinc := list.ToZinc()
	if listZinc != "[5.5, 23:07:10, @null]" {
		t.Error(listZinc)
	}
}
