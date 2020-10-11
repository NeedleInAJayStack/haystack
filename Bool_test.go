package haystack

import "testing"

func TestBool_ToZinc(t *testing.T) {
	valTest_ToZinc(TRUE, "T", t)
	valTest_ToZinc(FALSE, "F", t)
}

func TestBool_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(TRUE, "true", t)
	valTest_MarshalJSON(FALSE, "false", t)
}
