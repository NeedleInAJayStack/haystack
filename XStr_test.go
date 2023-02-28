package haystack

import (
	"encoding/json"
	"testing"
)

func TestXStr_ToZinc(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_ToZinc(xStr, "Str(\"hello world\")", t)
}

func TestXStr_MarhsalJSON(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_MarshalJSON(xStr, "\"x:Str:hello world\"", t)
}

func TestXStr_UnmarshalJSON(t *testing.T) {
	var val XStr
	json.Unmarshal([]byte("\"x:Str:hello world\""), &val)
	valTest_ToZinc(val, "Str(\"hello world\")", t)
}

func TestXStr_MarhsalHAYSON(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_MarshalHayson(xStr, "{\"_kind\":\"xstr\",\"type\":\"Str\",\"val\":\"hello world\"}", t)
}
