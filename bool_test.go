package haystack

import "testing"

func TestBool_ToZinc(t *testing.T) {
	trueBool := NewBool(true)
	trueStr := trueBool.ToZinc()
	if trueStr != "T" {
		t.Error(trueStr)
	}

	falseBool := NewBool(false)
	falseStr := falseBool.ToZinc()
	if falseStr != "F" {
		t.Error(falseStr)
	}
}
