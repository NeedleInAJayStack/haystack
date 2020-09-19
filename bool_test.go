package haystack

import "testing"

func TestBool_ToZinc(t *testing.T) {
	trueBool := Bool{val: true}
	trueStr := trueBool.ToZinc()
	if trueStr != "T" {
		t.Error(trueStr)
	}

	falseBool := Bool{val: false}
	falseStr := falseBool.ToZinc()
	if falseStr != "F" {
		t.Error(falseStr)
	}
}
