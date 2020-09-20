package haystack

import (
	"strings"
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
