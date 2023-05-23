package haystack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRef_ToZinc(t *testing.T) {
	assert.Equal(t, NewRef("123-abc", "").ToZinc(), "@123-abc")
	assert.Equal(t, NewRef("123-abc", "Name").ToZinc(), "@123-abc \"Name\"")
}

func TestRef_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewRef("123-abc", ""), "\"r:123-abc\"", t)
	valTest_MarshalJSON(NewRef("123-abc", "Name"), "\"r:123-abc Name\"", t)
}

func TestRef_UnmarshalJSON(t *testing.T) {
	var refNoDis Ref
	json.Unmarshal([]byte("\"r:123-abc\""), &refNoDis)
	assert.Equal(t, refNoDis, NewRef("123-abc", ""))

	var refDis Ref
	json.Unmarshal([]byte("\"r:123-abc Name\""), &refDis)
	assert.Equal(t, refDis, NewRef("123-abc", "Name"))
}

func TestRef_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewRef("123-abc", ""), "{\"_kind\":\"ref\",\"val\":\"123-abc\"}", t)
	valTest_MarshalHayson(NewRef("123-abc", "Name"), "{\"_kind\":\"ref\",\"val\":\"123-abc\",\"dis\":\"Name\"}", t)
}
