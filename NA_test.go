package haystack

import (
	"encoding/json"
	"testing"
)

func TestNA_ToZinc(t *testing.T) {
	valTest_ToZinc(NewNA(), "NA", t)
}

func TestNA_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNA(), "\"z:\"", t)
}

func TestNA_UnmarshalJSON(t *testing.T) {
	var na NA
	json.Unmarshal([]byte("\"z:\""), &na)
	valTest_ToZinc(na, "NA", t)
}

func TestNA_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNA(), "{\"_kind\":\"na\"}", t)
}
