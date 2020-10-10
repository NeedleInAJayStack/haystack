package haystack

import "testing"

func TestNull_ToZinc(t *testing.T) {
	null := NewNull()
	nullStr := null.ToZinc()
	if nullStr != "N" {
		t.Error(nullStr)
	}
}

func TestNull_ToJSON(t *testing.T) {
	null := NewNull()
	nullStr := null.ToJSON()
	if nullStr != "null" {
		t.Error(nullStr)
	}
}
