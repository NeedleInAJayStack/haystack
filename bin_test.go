package haystack

import "testing"

func TestBin_ToZinc(t *testing.T) {
	bin := NewBin("text/plain")
	binStr := bin.ToZinc()
	if binStr != "Bin(\"text/plain\")" {
		t.Error(binStr)
	}
}
