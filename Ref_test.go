package haystack

import (
	"testing"
)

func TestRef_ToZinc(t *testing.T) {
	valTest_ToZinc(NewRef("123-abc", ""), "@123-abc", t)
	valTest_ToZinc(NewRef("123-abc", "Name"), "@123-abc \"Name\"", t)
}

func TestRef_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewRef("123-abc", ""), "\"r:123-abc\"", t)
	valTest_MarshalJSON(NewRef("123-abc", "Name"), "\"r:123-abc Name\"", t)
}

func TestRef_UnmarshalJSON(t *testing.T) {
	var refNoDis Ref
	valTest_UnmarshalJSON("\"r:123-abc\"", refNoDis, "@123-abc", t)
	var refDis Ref
	valTest_UnmarshalJSON("\"r:123-abc Name\"", refNoDis, "@123-abc \"Name\"", t)
}

func TestRef_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewRef("123-abc", ""), "{\"_kind\":\"ref\",\"val\":\"123-abc\"}", t)
	valTest_MarshalHayson(NewRef("123-abc", "Name"), "{\"_kind\":\"ref\",\"val\":\"123-abc\",\"dis\":\"Name\"}", t)
}
