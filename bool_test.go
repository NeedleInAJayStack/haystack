package haystack

import "testing"

func TestBool_ToZinc(t *testing.T) {
	trueBool := TRUE
	trueStr := trueBool.ToZinc()
	if trueStr != "T" {
		t.Error(trueStr)
	}

	falseBool := FALSE
	falseStr := falseBool.ToZinc()
	if falseStr != "F" {
		t.Error(falseStr)
	}
}
