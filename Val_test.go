package haystack

import "testing"

func valTest_ToZinc(val Val, zinc string, t *testing.T) {
	valStr := val.ToZinc()
	if valStr != zinc {
		t.Error(valStr + " != " + zinc)
	}
}

func valTest_ToZinc_Grid(val Val, zinc string, t *testing.T) {
	// Customization to provide easier formatting for multiline grid output
	valStr := val.ToZinc()
	if valStr != zinc {
		t.Error("\nACTUAL:\n" + valStr + "\n\nEXPECT:\n" + zinc)
	}
}

func valTest_MarshalJSON(val Val, json string, t *testing.T) {
	valBytes, marshalErr := val.MarshalJSON()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	valStr := string(valBytes)
	if valStr != json {
		t.Error(valStr + " != " + json)
	}
}

func valTest_MarshalHayson(val HaysonMarshaller, hayson string, t *testing.T) {
	valBytes, marshalErr := val.MarshalHayson()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	valStr := string(valBytes)
	if valStr != hayson {
		t.Error(valStr + " != " + hayson)
	}
}
