package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNA_ToZinc(t *testing.T) {
	assert.Equal(t, NewNA().ToZinc(), "NA")
}

func TestNA_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNA(), "\"z:\"", t)
}

func TestNA_UnmarshalJSON(t *testing.T) {
	var na NA
	json.Unmarshal([]byte("\"z:\""), &na)
	assert.Equal(t, na, NewNA())
}

func TestNA_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNA(), "{\"_kind\":\"na\"}", t)
}
