package haystack

import (
	"encoding/json"
	"testing"
)

func TestUri_ToZinc(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_ToZinc(uri, "`http://www.project-haystack.org`", t)
}

func TestUri_MarshalJSON(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_MarshalJSON(uri, "\"u:http://www.project-haystack.org\"", t)
}

func TestUri_UnmarshalJSON(t *testing.T) {
	jsonStr := "\"u:http://www.project-haystack.org\""

	var val Uri
	err := json.Unmarshal([]byte(jsonStr), &val)
	if err != nil {
		t.Error(err)
	}
	valStr := val.ToZinc()
	if valStr != "`http://www.project-haystack.org`" {
		t.Error(valStr + " != " + "`http://www.project-haystack.org`")
	}
}

func TestUri_MarshalHayson(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_MarshalHayson(uri, "{\"_kind\":\"uri\",\"val\":\"http://www.project-haystack.org\"}", t)
}
