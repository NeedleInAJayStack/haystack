package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbol_ToZinc(t *testing.T) {
	assert.Equal(t, NewSymbol("foo").ToZinc(), "^foo")
}

func TestSymbol_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewSymbol("foo"), "\"y:foo\"", t)
}

func TestSymbol_UnmarshalJSON(t *testing.T) {
	var val Symbol
	json.Unmarshal([]byte("\"y:foo\""), &val)
	assert.Equal(t, val, NewSymbol("foo"))
}

func TestSymbol_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewSymbol("foo"), "{\"_kind\":\"symbol\",\"val\":\"foo\"}", t)
}
