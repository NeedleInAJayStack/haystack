package haystack

import "testing"

func TestNA_ToZinc(t *testing.T) {
	na := NewNA()
	valTest_ToZinc(na, "NA", t)
}

func TestNA_MarshalJSON(t *testing.T) {
	na := NewNA()
	valTest_MarshalJSON(na, "\"z:\"", t)
}

func TestNA_MarshalHayson(t *testing.T) {
	na := NewNA()
	valTest_MarshalHayson(na, "{\"_kind\":\"na\"}", t)
}
