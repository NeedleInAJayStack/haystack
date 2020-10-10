package haystack

import (
	"testing"
)

func TestRef_ToZinc(t *testing.T) {
	refNoDis := NewRef("123-abc", "")
	refNoDisStr := refNoDis.ToZinc()
	if refNoDisStr != "@123-abc" {
		t.Error(refNoDisStr)
	}

	refDis := NewRef("123-abc", "Name")
	refDisStr := refDis.ToZinc()
	if refDisStr != "@123-abc \"Name\"" {
		t.Error(refDisStr)
	}
}

func TestRef_ToJSON(t *testing.T) {
	refNoDis := NewRef("123-abc", "")
	refNoDisStr := refNoDis.ToJSON()
	if refNoDisStr != "r:123-abc" {
		t.Error(refNoDisStr)
	}

	refDis := NewRef("123-abc", "Name")
	refDisStr := refDis.ToJSON()
	if refDisStr != "r:123-abc Name" {
		t.Error(refDisStr)
	}
}
