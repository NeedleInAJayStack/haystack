package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXStr_ToZinc(t *testing.T) {
	assert.Equal(t, NewXStr("Str", "hello world").ToZinc(), "Str(\"hello world\")")
}

func TestXStr_MarhsalJSON(t *testing.T) {
	valTest_MarshalJSON(NewXStr("Str", "hello world"), "\"x:Str:hello world\"", t)
}

func TestXStr_UnmarshalJSON(t *testing.T) {
	var val XStr
	json.Unmarshal([]byte("\"x:Str:hello world\""), &val)
	assert.Equal(t, val, NewXStr("Str", "hello world"))
}

func TestXStr_MarhsalHAYSON(t *testing.T) {
	valTest_MarshalHayson(
		NewXStr("Str", "hello world"),
		"{\"_kind\":\"xstr\",\"type\":\"Str\",\"val\":\"hello world\"}",
		t,
	)
}
