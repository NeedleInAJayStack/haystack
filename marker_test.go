package haystack

import "testing"

func TestMarker_ToZinc(t *testing.T) {
	marker := NewMarker()
	markerZinc := marker.ToZinc()
	if markerZinc != "M" {
		t.Error(markerZinc)
	}
}
