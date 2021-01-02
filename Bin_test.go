package haystack

import (
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
	val := NewBin("")
	valTest_UnmarshalJSON("\"b:text/plain\"", val, "Bin(\"text/plain\")", t)
}

func TestBin_MarshalHayson(t *testing.T) {
	bin := NewBin("text/plain")
	valTest_MarshalHayson(bin, "{\"_kind\":\"bin\",\"mime\":\"text/plain\"}", t)
}
