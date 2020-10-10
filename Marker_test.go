package haystack

import "testing"

func TestMarker_ToZinc(t *testing.T) {
	marker := NewMarker()
	markerZinc := marker.ToZinc()
	if markerZinc != "M" {
		t.Error(markerZinc)
	}
}

func TestMarker_ToJSON(t *testing.T) {
	marker := NewMarker()
	markerJSON := marker.ToJSON()
	if markerJSON != "m:" {
		t.Error(markerJSON)
	}
}

func TestRemove_ToZinc(t *testing.T) {
	remove := NewRemove()
	removeZinc := remove.ToZinc()
	if removeZinc != "R" {
		t.Error(removeZinc)
	}
}

func TestRemove_ToJSON(t *testing.T) {
	remove := NewRemove()
	removeJSON := remove.ToJSON()
	if removeJSON != "x:" {
		t.Error(removeJSON)
	}
}
