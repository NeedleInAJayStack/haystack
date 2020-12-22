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

func TestNull_MarshalHAYSON(t *testing.T) {
	null := NewNull()
	valTest_MarshalHAYSON(null, "null", t)
}
