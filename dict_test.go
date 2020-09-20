package haystack

import (
	"testing"
)

func TestDict_ToZinc(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	dictZinc := dict.ToZinc()
	if dictZinc != "{area:35000ft² dis:\"Building\" site}" {
		t.Error(dictZinc)
	}
}
