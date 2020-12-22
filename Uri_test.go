package haystack

import "testing"

func TestUri_ToZinc(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_ToZinc(uri, "`http://www.project-haystack.org`", t)
}

func TestUri_MarshalJSON(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_MarshalJSON(uri, "\"u:http://www.project-haystack.org\"", t)
}

func TestUri_MarshalHAYSON(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	valTest_MarshalHAYSON(uri, "{\"_kind\":\"uri\",\"val\":\"http://www.project-haystack.org\"}", t)
}
