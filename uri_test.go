package haystack

import "testing"

func TestUri_toZinc(t *testing.T) {
	uri := Uri{val: "http://www.project-haystack.org"}
	got := uri.toZinc()
	if got != "`http://www.project-haystack.org`" {
		t.Error(got)
	}
}
