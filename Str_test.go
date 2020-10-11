package haystack

import "testing"

func TestStr_ToZinc(t *testing.T) {
	easy := NewStr("hello world")
	valTest_ToZinc(easy, "\"hello world\"", t)

	hard := NewStr("this 1s A #more \n complex \\one")
	valTest_ToZinc(hard, "\"this 1s A #more \\n complex \\\\one\"", t)
}

func TestStr_MarshalJSON(t *testing.T) {
	easy := NewStr("hello world")
	valTest_MarshalJSON(easy, "\"hello world\"", t)

	hasColon := NewStr("https://project-haystack.org/")
	valTest_MarshalJSON(hasColon, "\"s:https://project-haystack.org/\"", t)
}
