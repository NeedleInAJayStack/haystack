package haystack

import "testing"

func TestMarker_ToZinc(t *testing.T) {
	marker := NewMarker()
	valTest_ToZinc(marker, "M", t)
}

func TestMarker_MarshalJSON(t *testing.T) {
	marker := NewMarker()
	valTest_MarshalJSON(marker, "\"m:\"", t)
}

func TestMarker_MarshalHayson(t *testing.T) {
	marker := NewMarker()
	valTest_MarshalHayson(marker, "{\"_kind\":\"marker\"}", t)
}

func TestRemove_ToZinc(t *testing.T) {
	remove := NewRemove()
	valTest_ToZinc(remove, "R", t)
}

func TestRemove_MarshalJSON(t *testing.T) {
	remove := NewRemove()
	valTest_MarshalJSON(remove, "\"-:\"", t)
}

func TestRemove_MarshalHayson(t *testing.T) {
	remove := NewRemove()
	valTest_MarshalHayson(remove, "{\"_kind\":\"remove\"}", t)
}
