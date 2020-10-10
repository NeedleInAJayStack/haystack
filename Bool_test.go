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

func TestBool_ToJSON(t *testing.T) {
	trueBool := TRUE
	trueStr := trueBool.ToJSON()
	if trueStr != "true" {
		t.Error(trueStr)
	}

	falseBool := FALSE
	falseStr := falseBool.ToJSON()
	if falseStr != "false" {
		t.Error(falseStr)
	}
}
