package haystack

import (
	"testing"
)

func TestRef_toZinc(t *testing.T) {
	refNoDis := Ref{val: "123-abc"}
	refNoDisZinc := refNoDis.toZinc()
	if refNoDisZinc != "@123-abc" {
		t.Error(refNoDisZinc)
	}

	refDis := Ref{val: "123-abc", dis: "Name"}
	refDisZinc := refDis.toZinc()
	if refDisZinc != "@123-abc \"Name\"" {
		t.Error(refDisZinc)
	}
}
