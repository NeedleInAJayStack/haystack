package haystack

import "testing"

func TestUri_toZinc(t *testing.T) {
	uri := Uri{val: "http://www.projecthaystack.org"}
	got := uri.toZinc()
	if got != "`http://www.projecthaystack.org`" {
		t.Error(got)
	}
}
