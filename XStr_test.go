package haystack

import "testing"

func TestXStr_ToZinc(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_ToZinc(xStr, "Str(\"hello world\")", t)
}

func TestXStr_MarhsalJSON(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_MarshalJSON(xStr, "\"x:Str:hello world\"", t)
}
