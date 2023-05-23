package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBin_ToZinc(t *testing.T) {
	assert.Equal(t, NewBin("text/plain").ToZinc(), "Bin(\"text/plain\")")
}

func TestBin_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewBin("text/plain"), "\"b:text/plain\"", t)
}

func TestBin_UnmarshalJSON(t *testing.T) {
	var bin Bin
	json.Unmarshal([]byte("\"b:text/plain\""), &bin)
	assert.Equal(t, bin, NewBin("text/plain"))
}

func TestBin_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewBin("text/plain"), "{\"_kind\":\"bin\",\"mime\":\"text/plain\"}", t)
}
