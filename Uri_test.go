package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUri_ToZinc(t *testing.T) {
	assert.Equal(t, NewUri("http://www.project-haystack.org").ToZinc(), "`http://www.project-haystack.org`")
}

func TestUri_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewUri("http://www.project-haystack.org"), "\"u:http://www.project-haystack.org\"", t)
}

func TestUri_UnmarshalJSON(t *testing.T) {
	var val Uri
	json.Unmarshal([]byte("\"u:http://www.project-haystack.org\""), &val)
	assert.Equal(t, val, NewUri("http://www.project-haystack.org"))
}

func TestUri_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(
		NewUri("http://www.project-haystack.org"),
		"{\"_kind\":\"uri\",\"val\":\"http://www.project-haystack.org\"}",
		t,
	)
}
