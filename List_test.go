package haystack

import (
	"encoding/json"
	"testing"
)

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

func TestList_UnmarshalJSON(t *testing.T) {
	var list List
	json.Unmarshal([]byte("[\"n:5.5\",\"h:23:07:10\",\"r:null\"]"), &list)
	valTest_ToZinc(list, "[5.5, 23:07:10, @null]", t)
}

func TestList_MarshalHayson(t *testing.T) {
	list := NewList(
		[]Val{
			NewNumber(5.5, ""),
			NewTime(23, 7, 10, 0),
			NewRef("null", ""),
		},
	)
	valTest_MarshalHayson(list, "[{\"_kind\":\"number\",\"val\":5.5},{\"_kind\":\"time\",\"val\":\"23:07:10\"},{\"_kind\":\"ref\",\"val\":\"null\"}]", t)
}
