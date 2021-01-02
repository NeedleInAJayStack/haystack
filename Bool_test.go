package haystack

import (
	"testing"
)

func TestBool_ToZinc(t *testing.T) {
	valTest_ToZinc(NewBool(true), "T", t)
	valTest_ToZinc(NewBool(false), "F", t)
}

func TestBool_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewBool(true), "true", t)
	valTest_MarshalJSON(NewBool(false), "false", t)
}

func TestBool_UnmarshalJSON(t *testing.T) {
	val := NewBool(false)
	valTest_UnmarshalJSON("true", val, "T", t)
}

func TestBool_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewBool(true), "true", t)
	valTest_MarshalHayson(NewBool(false), "false", t)
}
