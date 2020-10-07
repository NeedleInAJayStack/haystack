package haystack

import "testing"

func TestNA_ToZinc(t *testing.T) {
	na := NewNA()
	naZinc := na.ToZinc()
	if naZinc != "NA" {
		t.Error(naZinc)
	}
}
