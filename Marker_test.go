package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarker_ToZinc(t *testing.T) {
	assert.Equal(t, NewMarker().ToZinc(), "M")
}

func TestMarker_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewMarker(), "\"m:\"", t)
}

func TestMarker_UnmarshalJSON(t *testing.T) {
	var marker Marker
	json.Unmarshal([]byte("\"m:\""), &marker)
	assert.Equal(t, marker, NewMarker())
}

func TestMarker_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewMarker(), "{\"_kind\":\"marker\"}", t)
}

func TestRemove_ToZinc(t *testing.T) {
	assert.Equal(t, NewRemove().ToZinc(), "R")
}

func TestRemove_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewRemove(), "\"-:\"", t)
}

func TestRemove_UnmarshalJSON(t *testing.T) {
	var remove Remove
	json.Unmarshal([]byte("\"-:\""), &remove)
	assert.Equal(t, remove, NewRemove())
}

func TestRemove_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewRemove(), "{\"_kind\":\"remove\"}", t)
}
