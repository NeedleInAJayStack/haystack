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
	valTest_ToZinc(NewCoord(41.534, 111.478), "C(41.534,111.478)", t)
}

func TestCoord_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewCoord(41.534, 111.478), "\"c:41.534,111.478\"", t)
}

func TestCoord_UnmarshalJSON(t *testing.T) {
	val := NewCoord(0, 0)
	valTest_UnmarshalJSON("\"c:41.534,111.478\"", val, "C(41.534,111.478)", t)
}

func TestCoord_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewCoord(41.534, 111.478), "{\"_kind\":\"coord\",\"lat\":41.534,\"lng\":111.478}", t)
}
