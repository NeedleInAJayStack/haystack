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

func TestRemove_ToZinc(t *testing.T) {
	remove := NewRemove()
	valTest_ToZinc(remove, "R", t)
}

func TestRemove_MarshalJSON(t *testing.T) {
	remove := NewRemove()
	valTest_MarshalJSON(remove, "\"-:\"", t)
}
