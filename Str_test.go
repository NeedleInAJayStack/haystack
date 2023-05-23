package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStr_ToZinc(t *testing.T) {
	assert.Equal(t, NewStr("hello world").ToZinc(), "\"hello world\"")
	assert.Equal(t, NewStr("this 1s A #more \n complex \\one").ToZinc(), "\"this 1s A #more \\n complex \\\\one\"")
}

func TestStr_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewStr("hello world"), "\"hello world\"", t)
	valTest_MarshalJSON(NewStr("https://project-haystack.org/"), "\"s:https://project-haystack.org/\"", t)
}

func TestStr_UnmarshalJSON(t *testing.T) {
	var val Str
	json.Unmarshal([]byte("\"s:https://project-haystack.org/\""), &val)
	assert.Equal(t, val, NewStr("https://project-haystack.org/"))
}

func TestStr_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewStr("hello world"), "\"hello world\"", t)
	valTest_MarshalHayson(NewStr("https://project-haystack.org/"), "\"https://project-haystack.org/\"", t)
}
