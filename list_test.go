package haystack

import "testing"

func TestList_ToZinc(t *testing.T) {
	list := List{
		vals: []Val{
			&Number{val: 5.5},
			&Time{hour: 23, min: 7, sec: 10},
			&Ref{id: "null"},
		},
	}
	listZinc := list.ToZinc()
	if listZinc != "[5.5, 23:07:10, @null]" {
		t.Error(listZinc)
	}
}
