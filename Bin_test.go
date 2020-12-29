package haystack

import (
	"encoding/json"
	"testing"
)

func TestBin_ToZinc(t *testing.T) {
	bin := NewBin("text/plain")
	valTest_ToZinc(bin, "Bin(\"text/plain\")", t)
}

func TestBin_MarshalJSON(t *testing.T) {
	bin := NewBin("text/plain")
	valTest_MarshalJSON(bin, "\"b:text/plain\"", t)
}

func TestBin_UnmarshalJSON(t *testing.T) {
	jsonStr := "\"b:text/plain\""

	var val Bin
	err := json.Unmarshal([]byte(jsonStr), &val)
	if err != nil {
		t.Error(err)
	}
	valStr := val.ToZinc()
	if valStr != "Bin(\"text/plain\")" {
		t.Error(valStr + " != " + "Bin(\"text/plain\")")
	}
}

func TestBin_MarshalHayson(t *testing.T) {
	bin := NewBin("text/plain")
	valTest_MarshalHayson(bin, "{\"_kind\":\"bin\",\"mime\":\"text/plain\"}", t)
}
