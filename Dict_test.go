package haystack

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDict_Names(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)

	names := dict.Names()
	sort.Strings(names)
	if names[0] != "area" {
		t.Error("Names missing 'area' field")
	}
	if names[1] != "dis" {
		t.Error("Names missing 'dis' field")
	}
	if names[2] != "site" {
		t.Error("Names missing 'site' field")
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

	if newDict.Get("geoState").ToZinc() != NewStr("UT").ToZinc() {
		t.Error("Dict.Set didn't set a new field correctly")
	}

	// Ensure original wasn't changed
	if dict.Get("geoState").ToZinc() != NewNull().ToZinc() {
		t.Error("Dict.Set changed the state of the original Dict")
	}

	overrideDict := dict.Set("dis", NewStr("Different Building"))
	if overrideDict.Get("dis").ToZinc() != NewStr("Different Building").ToZinc() {
		t.Error("Dict.Set didn't override an existing value correctly")
	}
}

func TestDict_SetAll(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	newDict := dict.SetAll(map[string]Val{
		"geoState": NewStr("UT"),
		"geoCity":  NewStr("Salt Lake City"),
	})

	assert.Equal(t, newDict.Get("geoState"), NewStr("UT"), "Dict.SetAll didn't set all values correctly")
	assert.Equal(t, newDict.Get("geoCity"), NewStr("Salt Lake City"), "Dict.SetAll didn't set all values correctly")

	// Ensure original wasn't changed
	assert.Equal(t, dict.Get("geoState"), NewNull(), "Dict.SetAll changed the state of the original Dict")
}

func TestDict_Size(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	assert.Equal(t, dict.Size(), 3, "Dict.Size returned an incorrect value")
}

func TestDict_IsEmpty(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	assert.False(t, dict.IsEmpty(), "Dict.IsEmpty returned true on a non-empty grid")

	emptyDict := NewDict(
		map[string]Val{},
	)
	assert.True(t, emptyDict.IsEmpty(), "Dict.IsEmpty returned false on an empty grid")
}

func TestDict_ToZinc(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	assert.Equal(t, dict.ToZinc(), "{area:35000ft² dis:\"Building\" site}")
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

func TestDict_UnmarshalJSON(t *testing.T) {
	var dict Dict
	json.Unmarshal([]byte("{\"area\":\"n:35000 ft²\",\"dis\":\"Building\",\"site\":\"m:\"}"), &dict)
	assert.Equal(t, dict.ToZinc(), "{area:35000ft² dis:\"Building\" site}")
}

func TestDict_MarshalHayson(t *testing.T) {
	dict := NewDict(
		map[string]Val{
			"dis":  NewStr("Building"),
			"site": NewMarker(),
			"area": NewNumber(35000.0, "ft²"),
		},
	)
	valTest_MarshalHayson(dict, "{\"_kind\":\"dict\",\"area\":{\"_kind\":\"number\",\"val\":35000,\"unit\":\"ft²\"},\"dis\":\"Building\",\"site\":{\"_kind\":\"marker\"}}", t)
}
