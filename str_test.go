package haystack

import "testing"

func TestStr_toZinc(t *testing.T) {
	easy := Str{val: "hello world"}
	easyStr := easy.toZinc()
	if easyStr != "\"hello world\"" {
		t.Error(easyStr)
	}

	hard := Str{val: "this 1s A #more \n complex \\one"}
	hardStr := hard.toZinc()
	if hardStr != "\"this 1s A #more \\n complex \\\\one\"" {
		t.Error(hardStr)
	}
}
