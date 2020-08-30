package haystack

import "testing"

func TestDict_toZinc(t *testing.T) {
	dict := Dict{
		items: map[string]Val{
			"dis":  &Str{val: "Building"},
			"site": &Marker{},
			"area": &Number{val: 35000.0, unit: "ft²"},
		},
	}
	dictZinc := dict.toZinc()
	// Might fail if order is not guaranteed
	if dictZinc != "{dis:\"Building\" site area:35000ft²}" {
		t.Error(dictZinc)
	}
}
