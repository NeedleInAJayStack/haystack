package haystack

import "testing"

func TestDict_toZinc(t *testing.T) {
	dict := Dict{
		items: map[string]Val{
			"area":        &Number{val: 5.5},
			"currentTime": &Time{hour: 23, min: 7, sec: 10},
			"id":          &Ref{val: "null"},
		},
	}
	dictZinc := dict.toZinc()
	if dictZinc != "{area:5.5,currentTime:23:07:10,id:@null}" {
		t.Error(dictZinc)
	}
}
