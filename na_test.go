package haystack

import "testing"

func TestNA_toZinc(t *testing.T) {
	na := NA{}
	naZinc := na.toZinc()
	if naZinc != "NA" {
		t.Error(naZinc)
	}
}
