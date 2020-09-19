package haystack

import "testing"

func TestNA_ToZinc(t *testing.T) {
	na := NA{}
	naZinc := na.ToZinc()
	if naZinc != "NA" {
		t.Error(naZinc)
	}
}
