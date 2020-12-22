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

func TestMarker_MarshalHAYSON(t *testing.T) {
	marker := NewMarker()
	valTest_MarshalHAYSON(marker, "{\"_kind\":\"marker\"}", t)
}

func TestRemove_ToZinc(t *testing.T) {
	remove := NewRemove()
	valTest_ToZinc(remove, "R", t)
}

func TestRemove_MarshalJSON(t *testing.T) {
	remove := NewRemove()
	valTest_MarshalJSON(remove, "\"-:\"", t)
}

func TestRemove_MarshalHAYSON(t *testing.T) {
	remove := NewRemove()
	valTest_MarshalHAYSON(remove, "{\"_kind\":\"remove\"}", t)
}
