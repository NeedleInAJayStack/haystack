package haystack

import "testing"

func TestNA_ToZinc(t *testing.T) {
	na := NewNA()
	naStr := na.ToZinc()
	if naStr != "NA" {
		t.Error(naStr)
	}
}

func TestNA_ToJSON(t *testing.T) {
	na := NewNA()
	naStr := na.ToJSON()
	if naStr != "z:" {
		t.Error(naStr)
	}
}
