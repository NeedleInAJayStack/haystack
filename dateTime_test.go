package haystack

import "testing"

func TestDateTime_toZinc(t *testing.T) {
	date := Date{year: 2020, month: 8, day: 17}
	time := Time{hour: 23, min: 7, sec: 10}
	dateTime := DateTime{date: date, time: time}
	dateTimeZinc := dateTime.toZinc()
	if dateTimeZinc != "2020-08-17T23:07:10" {
		t.Error(dateTimeZinc)
	}
}
