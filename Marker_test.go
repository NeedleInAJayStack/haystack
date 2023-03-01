package haystack

import (
	"encoding/json"
	"testing"
)

func TestMarker_ToZinc(t *testing.T) {
	valTest_ToZinc(NewMarker(), "M", t)
}

func TestMarker_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewMarker(), "\"m:\"", t)
}

func TestMarker_UnmarshalJSON(t *testing.T) {
	var marker Marker
	json.Unmarshal([]byte("\"m:\""), &marker)
	valTest_ToZinc(marker, "M", t)
}

func TestMarker_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewMarker(), "{\"_kind\":\"marker\"}", t)
}

func TestRemove_ToZinc(t *testing.T) {
	valTest_ToZinc(NewRemove(), "R", t)
}

func TestRemove_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewRemove(), "\"-:\"", t)
}

func TestRemove_UnmarshalJSON(t *testing.T) {
	var remove Remove
	json.Unmarshal([]byte("\"-:\""), &remove)
	valTest_ToZinc(remove, "R", t)
}

func TestRemove_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewRemove(), "{\"_kind\":\"remove\"}", t)
}
