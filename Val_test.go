package haystack

import (
	"encoding/json"
	"testing"
)

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

func valTest_UnmarshalJSON(input string, val Val, expectedZinc string, t *testing.T) {
	unmarshalErr := json.Unmarshal([]byte(input), &val)
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
	}
	actualZinc := Val(val).ToZinc()
	if actualZinc != expectedZinc {
		t.Error(actualZinc + " != " + expectedZinc)
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
