package haystack

import (
	"encoding/json"
	"testing"
)

func TestSymbol_ToZinc(t *testing.T) {
	valTest_ToZinc(NewSymbol("foo"), "^foo", t)
}

func TestSymbol_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewSymbol("foo"), "\"y:foo\"", t)
}

func TestSymbol_UnmarshalJSON(t *testing.T) {
	var val Symbol
	json.Unmarshal([]byte("\"y:foo\""), &val)
	valTest_ToZinc(val, "^foo", t)
}

func TestSymbol_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewSymbol("foo"), "{\"_kind\":\"symbol\",\"val\":\"foo\"}", t)
}
