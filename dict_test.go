package haystack

import (
	"strings"
	"testing"
)

func TestDict_toZinc(t *testing.T) {
	dict := Dict{
		items: map[string]Val{
			"dis":  &Str{val: "Building"},
			"site": &Marker{},
			"area": &Number{val: 35000.0, unit: "ft²"},
		},
	}
	dictZinc := dict.toZinc()
	if !strings.Contains(dictZinc, "dis:\"Building\"") {
		t.Error(dictZinc)
	}
	if !strings.Contains(dictZinc, "site") {
		t.Error(dictZinc)
	}
	if !strings.Contains(dictZinc, "area:35000ft²") {
		t.Error(dictZinc)
	}
}
