package haystack

import (
	"encoding/json"
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
	var bool Bool
	json.Unmarshal([]byte("true"), &bool)
	valTest_ToZinc(bool, "T", t)
}

func TestBool_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewBool(true), "true", t)
	valTest_MarshalHayson(NewBool(false), "false", t)
}
