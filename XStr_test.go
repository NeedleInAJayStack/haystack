package haystack

import "testing"

func TestXStr_ToZinc(t *testing.T) {
	xStr := NewXStr("Str", "hello world")
	valTest_ToZinc(xStr, "Str(\"hello world\")", t)
}
