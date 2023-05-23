package haystack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func valTest_MarshalJSON(val Val, expected string, t *testing.T) {
	bytes, marshalErr := val.MarshalJSON()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	actual := string(bytes)
	assert.Equal(t, actual, expected)
}

func valTest_MarshalHayson(val Val, expected string, t *testing.T) {
	bytes, marshalErr := val.MarshalHayson()
	if marshalErr != nil {
		t.Error(marshalErr)
	}
	actual := string(bytes)
	assert.Equal(t, actual, expected)
}
