package haystack

import "testing"

func TestTime_equals(t *testing.T) {
	time1 := Time{hour: 23, min: 7, sec: 10, ms: 957}
	time2 := Time{hour: 23, min: 7, sec: 10, ms: 957}
	time3 := Time{hour: 0, min: 0, sec: 0, ms: 0}

	if !time1.equals(&time1) {
		t.Error("The same object doesn't equal itself")
	}
	if !time1.equals(&time2) {
		t.Error("Equivalent objects doesn't equal itself")
	}
	if !time2.equals(&time1) {
		t.Error("Ordering matters")
	}
	if time1.equals(&time3) {
		t.Error("Non-equivalent objects are equal")
	}
}

func TestTime_timeFromStr(t *testing.T) {
	noMs := "23:07:10"
	expNoMs := Time{hour: 23, min: 7, sec: 10}
	timeNoMs, err := timeFromStr(noMs)
	if err != nil {
		t.Error(err)
	}
	if !expNoMs.equals(&timeNoMs) {
		t.Error(timeNoMs)
	}

	ms := "23:07:10.957"
	expMs := Time{hour: 23, min: 7, sec: 10, ms: 957}
	timeMs, err := timeFromStr(ms)
	if err != nil {
		t.Error(err)
	}
	if !expMs.equals(&timeMs) {
		t.Error(timeMs)
	}
}

func TestTime_toZinc(t *testing.T) {
	timeNoMs := Time{hour: 23, min: 7, sec: 10}
	timeNoMsZinc := timeNoMs.toZinc()
	if timeNoMsZinc != "23:07:10" {
		t.Error(timeNoMsZinc)
	}

	timeMs := Time{hour: 23, min: 7, sec: 10, ms: 957}
	timeMsZinc := timeMs.toZinc()
	if timeMsZinc != "23:07:10.957" {
		t.Error(timeMsZinc)
	}

	timeOnesMs := Time{hour: 23, min: 7, sec: 10, ms: 2}
	timeOnesMsZinc := timeOnesMs.toZinc()
	if timeOnesMsZinc != "23:07:10.002" {
		t.Error(timeOnesMsZinc)
	}

	timeTensMs := Time{hour: 23, min: 7, sec: 10, ms: 56}
	timeTensMsZinc := timeTensMs.toZinc()
	if timeTensMsZinc != "23:07:10.056" {
		t.Error(timeTensMsZinc)
	}
}
