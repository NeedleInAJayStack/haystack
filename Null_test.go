package haystack

import "testing"

func TestNull_ToZinc(t *testing.T) {
	null := NewNull()
	valTest_ToZinc(null, "N", t)
}

func TestNull_MarshalJSON(t *testing.T) {
	null := NewNull()
	valTest_MarshalJSON(null, "null", t)
}

func TestNull_UnmarshalJSON(t *testing.T) {
	var val Null
	valTest_UnmarshalJSON("null", val, "N", t)
}

func TestNull_MarshalHayson(t *testing.T) {
	null := NewNull()
	valTest_MarshalHayson(null, "null", t)
}
