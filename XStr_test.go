package haystack

import "testing"

func TestXStr_ToZinc(t *testing.T) {
	easy := NewXStr("Str", "hello world")
	easyStr := easy.ToZinc()
	if easyStr != "Str(\"hello world\")" {
		t.Error(easyStr)
	}
}
