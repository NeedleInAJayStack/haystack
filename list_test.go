package haystack

import "testing"

func TestList_toZinc(t *testing.T) {
	list := List{
		items: []Val{
			&Number{val: 5.5},
			&Time{hour: 23, min: 7, sec: 10},
			&Ref{val: "null"},
		},
	}
	listZinc := list.toZinc()
	if listZinc != "[5.5, 23:07:10, @null]" {
		t.Error(listZinc)
	}
}
