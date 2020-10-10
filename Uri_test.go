package haystack

import "testing"

func TestUri_ToZinc(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	uriStr := uri.ToZinc()
	if uriStr != "`http://www.project-haystack.org`" {
		t.Error(uriStr)
	}
}

func TestUri_ToJSON(t *testing.T) {
	uri := NewUri("http://www.project-haystack.org")
	uriStr := uri.ToJSON()
	if uriStr != "u:http://www.project-haystack.org" {
		t.Error(uriStr)
	}
}
