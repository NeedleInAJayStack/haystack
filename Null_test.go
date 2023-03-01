package haystack

import (
	"encoding/json"
	"testing"
)

func TestNull_ToZinc(t *testing.T) {
	valTest_ToZinc(NewNull(), "N", t)
}

func TestNull_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNull(), "null", t)
}

func TestNull_UnmarshalJSON(t *testing.T) {
	var null Null
	json.Unmarshal([]byte("null"), &null)
	valTest_ToZinc(null, "N", t)
}

func TestNull_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNull(), "null", t)
}
