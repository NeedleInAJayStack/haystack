package haystack

import "testing"

func TestUri_ToZinc(t *testing.T) {
	uri := Uri{val: "http://www.project-haystack.org"}
	got := uri.ToZinc()
	if got != "`http://www.project-haystack.org`" {
		t.Error(got)
	}
}
