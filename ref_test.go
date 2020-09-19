package haystack

import (
	"testing"
)

func TestRef_ToZinc(t *testing.T) {
	refNoDis := Ref{val: "123-abc"}
	refNoDisZinc := refNoDis.ToZinc()
	if refNoDisZinc != "@123-abc" {
		t.Error(refNoDisZinc)
	}

	refDis := Ref{val: "123-abc", dis: "Name"}
	refDisZinc := refDis.ToZinc()
	if refDisZinc != "@123-abc \"Name\"" {
		t.Error(refDisZinc)
	}
}
