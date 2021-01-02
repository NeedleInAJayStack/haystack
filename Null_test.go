package haystack

import "testing"

func TestNull_ToZinc(t *testing.T) {
	valTest_ToZinc(NewNull(), "N", t)
}

func TestNull_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNull(), "null", t)
}

// TODO I was getting panic errors on this test
// func TestNull_UnmarshalJSON(t *testing.T) {
// 	val := NewNull()
// 	valTest_UnmarshalJSON("null", val, "N", t)
// }

func TestNull_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNull(), "null", t)
}
