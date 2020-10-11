package haystack

import "testing"

func TestCoord_NewCoord(t *testing.T) {
	valid := NewCoord(41.534, 111.478)
	if valid.lat != 41.534 {
		t.Error(valid.lat)
	}
	if valid.lng != 111.478 {
		t.Error(valid.lng)
	}

	invalid := NewCoord(100.0, -50.0)
	if invalid.lat != 90.0 {
		t.Error(invalid.lat)
	}
	if invalid.lng != 0 {
		t.Error(invalid.lng)
	}
}

func TestCoord_ToZinc(t *testing.T) {
	coord := NewCoord(41.534, 111.478)
	valTest_ToZinc(coord, "C(41.534,111.478)", t)
}

func TestCoord_MarshalJSON(t *testing.T) {
	coord := NewCoord(41.534, 111.478)
	valTest_MarshalJSON(coord, "\"c:41.534,111.478\"", t)
}