package haystack

import "testing"

func TestList_Get(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	if list.Get(2).ToZinc() != NewRef("null", "").ToZinc() {
		t.Error("List.Get returned an incorrect value")
	}
}

func TestList_Size(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	if list.Size() != 3 {
		t.Error("List.Size returned an incorrect value")
	}
}

func TestList_MarshalJSON(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	valTest_MarshalJSON(list, "[\"n:5.5\",\"h:23:07:10\",\"r:null\"]", t)
}

func TestList_ToZinc(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	valTest_ToZinc(list, "[5.5, 23:07:10, @null]", t)
}
