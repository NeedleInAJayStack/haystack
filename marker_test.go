package haystack

import "testing"

func TestMarker_toZinc(t *testing.T) {
	marker := Marker{}
	markerZinc := marker.toZinc()
	if markerZinc != "M" {
		t.Error(markerZinc)
	}
}
