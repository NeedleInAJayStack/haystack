package haystack

import "testing"

func TestMarker_ToZinc(t *testing.T) {
	marker := Marker{}
	markerZinc := marker.ToZinc()
	if markerZinc != "M" {
		t.Error(markerZinc)
	}
}
