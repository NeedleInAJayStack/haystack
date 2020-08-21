package haystack

import "testing"

func TestXStr_toZinc(t *testing.T) {
	easy := XStr{valType: "Str", val: "hello world"}
	easyStr := easy.toZinc()
	if easyStr != "Str(\"hello world\")" {
		t.Error(easyStr)
	}
}
