package haystack

import "testing"

func TestDate_toZinc(t *testing.T) {
	date := Date{year: 2020, month: 8, day: 17}
	dateZinc := date.toZinc()
	if dateZinc != "2020-08-17" {
		t.Error(dateZinc)
	}
}
