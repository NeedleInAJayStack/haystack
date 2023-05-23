package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, NewCoord(41.534, 111.478).ToZinc(), "C(41.534,111.478)")
}

func TestCoord_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewCoord(41.534, 111.478), "\"c:41.534,111.478\"", t)
}

func TestCoord_UnmarshalJSON(t *testing.T) {
	var coord Coord
	json.Unmarshal([]byte("\"c:41.534,111.478\""), &coord)
	assert.Equal(t, coord, NewCoord(41.534, 111.478))
}

func TestCoord_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewCoord(41.534, 111.478), "{\"_kind\":\"coord\",\"lat\":41.534,\"lng\":111.478}", t)
}
