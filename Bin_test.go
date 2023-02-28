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
	var bin Bin
	json.Unmarshal([]byte("\"b:text/plain\""), &bin)
	valTest_ToZinc(bin, "Bin(\"text/plain\")", t)
}

func TestBin_MarshalHayson(t *testing.T) {
	bin := NewBin("text/plain")
	valTest_MarshalHayson(bin, "{\"_kind\":\"bin\",\"mime\":\"text/plain\"}", t)
}
