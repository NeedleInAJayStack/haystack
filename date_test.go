package haystack

import "testing"

func TestDate_equals(t *testing.T) {
	date1 := Date{year: 2020, month: 8, day: 17}
	date2 := Date{year: 2020, month: 8, day: 17}
	date3 := Date{year: 0, month: 0, day: 0}

	if !date1.equals(&date1) {
		t.Error("The same object doesn't equal itself")
	}
	if !date1.equals(&date2) {
		t.Error("Equivalent objects doesn't equal itself")
	}
	if !date2.equals(&date1) {
		t.Error("Ordering matters")
	}
	if date1.equals(&date3) {
		t.Error("Non-equivalent objects are equal")
	}
}

func TestDate_dateFromStr(t *testing.T) {
	dateStr := "2020-08-17"
	exp := Date{year: 2020, month: 8, day: 17}
	date, err := dateFromStr(dateStr)
	if err != nil {
		t.Error(err)
	}
	if !exp.equals(&date) {
		t.Error(date)
	}
}

func TestDate_toZinc(t *testing.T) {
	date := Date{year: 2020, month: 8, day: 17}
	dateZinc := date.toZinc()
	if dateZinc != "2020-08-17" {
		t.Error(dateZinc)
	}
}
