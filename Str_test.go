package haystack

import (
	"testing"
)

func TestStr_ToZinc(t *testing.T) {
	valTest_ToZinc(NewStr("hello world"), "\"hello world\"", t)
	valTest_ToZinc(NewStr("this 1s A #more \n complex \\one"), "\"this 1s A #more \\n complex \\\\one\"", t)
}

func TestStr_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewStr("hello world"), "\"hello world\"", t)
	valTest_MarshalJSON(NewStr("https://project-haystack.org/"), "\"s:https://project-haystack.org/\"", t)
}

func TestStr_UnmarshalJSON(t *testing.T) {
	var val Str
	valTest_UnmarshalJSON("\"s:https://project-haystack.org/\"", val, "\"https://project-haystack.org/\"", t)
}

func TestStr_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewStr("hello world"), "\"hello world\"", t)
	valTest_MarshalHayson(NewStr("https://project-haystack.org/"), "\"https://project-haystack.org/\"", t)
}
