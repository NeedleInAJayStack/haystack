package haystack

import "testing"

func TestNA_ToZinc(t *testing.T) {
	valTest_ToZinc(NewNA(), "NA", t)
}

func TestNA_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNA(), "\"z:\"", t)
}

func TestNA_UnmarshalJSON(t *testing.T) {
	val := NewNA()
	valTest_UnmarshalJSON("\"z:\"", val, "NA", t)
}

func TestNA_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNA(), "{\"_kind\":\"na\"}", t)
}
