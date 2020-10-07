package haystack

import "testing"

func TestStr_ToZinc(t *testing.T) {
	easy := NewStr("hello world")
	easyStr := easy.ToZinc()
	if easyStr != "\"hello world\"" {
		t.Error(easyStr)
	}

	hard := NewStr("this 1s A #more \n complex \\one")
	hardStr := hard.ToZinc()
	if hardStr != "\"this 1s A #more \\n complex \\\\one\"" {
		t.Error(hardStr)
	}
}
