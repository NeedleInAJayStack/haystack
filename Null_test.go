package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNull_ToZinc(t *testing.T) {
	assert.Equal(t, NewNull().ToZinc(), "N")
}

func TestNull_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewNull(), "null", t)
}

func TestNull_UnmarshalJSON(t *testing.T) {
	var null Null
	json.Unmarshal([]byte("null"), &null)
	assert.Equal(t, null, NewNull())
}

func TestNull_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewNull(), "null", t)
}
