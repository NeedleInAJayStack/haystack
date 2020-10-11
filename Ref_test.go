package haystack

import (
	"testing"
)

func TestRef_ToZinc(t *testing.T) {
	refNoDis := NewRef("123-abc", "")
	valTest_ToZinc(refNoDis, "@123-abc", t)

	refDis := NewRef("123-abc", "Name")
	valTest_ToZinc(refDis, "@123-abc \"Name\"", t)
}

func TestRef_MarshalJSON(t *testing.T) {
	refNoDis := NewRef("123-abc", "")
	valTest_MarshalJSON(refNoDis, "\"r:123-abc\"", t)

	refDis := NewRef("123-abc", "Name")
	valTest_MarshalJSON(refDis, "\"r:123-abc Name\"", t)
}
