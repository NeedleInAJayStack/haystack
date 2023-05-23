package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool_ToZinc(t *testing.T) {
	assert.Equal(t, NewBool(true).ToZinc(), "T")
	assert.Equal(t, NewBool(false).ToZinc(), "F")
}

func TestBool_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewBool(true), "true", t)
	valTest_MarshalJSON(NewBool(false), "false", t)
}

func TestBool_UnmarshalJSON(t *testing.T) {
	var bool Bool
	json.Unmarshal([]byte("true"), &bool)
	assert.Equal(t, bool, NewBool(true))
}

func TestBool_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewBool(true), "true", t)
	valTest_MarshalHayson(NewBool(false), "false", t)
}
