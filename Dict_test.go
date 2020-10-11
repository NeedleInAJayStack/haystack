package haystack

import (
	"fmt"
	"testing"
)

func TestDict_Get(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)

	dis := dict.Get("dis")
	if dis.ToZinc() != NewStr("Building").ToZinc() {
		t.Error(dis)
	}

	notHere := dict.Get("notHere")
	if notHere.ToZinc() != NewNull().ToZinc() {
		t.Error(notHere)
	}
}

func TestDict_Set(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	newDict := dict.Set("geoState", NewStr("UT"))

	geoState := newDict.Get("geoState")
	if geoState.ToZinc() != NewStr("UT").ToZinc() {
		t.Error("Dict.Set didn't set the value correctly")
	}

	// Ensure dict wasn't changed
	noGeoState := dict.Get("geoState")
	if noGeoState.ToZinc() != NewNull().ToZinc() {
		fmt.Println(noGeoState.ToZinc())
		t.Error("Dict.Set changed the state of the original Dict")
	}
}

func TestDict_ToZinc(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	valTest_ToZinc(dict, "{area:35000ft² dis:\"Building\" site}", t)
}

func TestDict_MarshalJSON(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	valTest_MarshalJSON(dict, "{\"area\":\"n:35000 ft²\",\"dis\":\"Building\",\"site\":\"m:\"}", t)
}
