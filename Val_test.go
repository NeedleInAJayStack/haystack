package haystack

import (
	"testing"
)

func valTest_Equal_Grid(actual Val, expected Val, t *testing.T) {
	// Compare based on zinc representation
	actualZinc := actual.ToZinc()
	expectedZinc := expected.ToZinc()
	if actualZinc != expectedZinc {
		t.Error("\nACTUAL:\n" + actualZinc + "\n\nEXPECT:\n" + expectedZinc)
	}
}

func valTest_ToZinc(val Val, expected string, t *testing.T) {
	actual := val.ToZinc()
	if actual != expected {
		t.Error(actual + " != " + expected)
	}
}

func valTest_ToZinc_Grid(val Val, expected string, t *testing.T) {
	// Customization to provide easier formatting for multiline grid output
	actual := val.ToZinc()
	if actual != expected {
		t.Error("\nACTUAL:\n" + actual + "\n\nEXPECT:\n" + expected)
	}
}

func valTest_MarshalJSON(val Val, expected string, t *testing.T) {
	bytes, marshalErr := val.MarshalJSON()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	actual := string(bytes)
	if actual != expected {
		t.Error(actual + " != " + expected)
	}
}

func valTest_MarshalHayson(val Val, expected string, t *testing.T) {
	bytes, marshalErr := val.MarshalHayson()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	actual := string(bytes)
	if actual != expected {
		t.Error(actual + " != " + expected)
	}
}
