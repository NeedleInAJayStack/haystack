package haystack

import (
	"encoding/json"
	"testing"
)

func valTest_Equal(actual Val, expected Val, t *testing.T) {
	// Compare based on zinc representation
	actualZinc := actual.ToZinc()
	expectedZinc := expected.ToZinc()
	if actualZinc != expectedZinc {
		t.Error(actualZinc + " != " + expectedZinc)
	}
}

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

func valTest_UnmarshalJSON(input string, val Val, expectedZinc string, t *testing.T) {
	unmarshalErr := json.Unmarshal([]byte(input), &val)
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
	}
	// We must add this because UnmarshalJSON([]byte("null")) is implemented as a no-op
	// meaning our Null unmarshaller doesn't work and we get nil.
	// See https://golang.org/pkg/encoding/json/#Unmarshal
	if val == nil {
		val = NewNull()
	}
	actualZinc := val.ToZinc()
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
