package haystack

import "testing"

func TestBool_toZinc(t *testing.T) {
	trueBool := Bool{val: true}
	trueStr := trueBool.toZinc()
	if trueStr != "T" {
		t.Error(trueStr)
	}

	falseBool := Bool{val: false}
	falseStr := falseBool.toZinc()
	if falseStr != "F" {
		t.Error(falseStr)
	}
}
