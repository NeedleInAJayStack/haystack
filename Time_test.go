package haystack

import "testing"

func TestTime_equals(t *testing.T) {
	time1 := NewTime(23, 7, 10, 957)
	time2 := NewTime(23, 7, 10, 957)
	time3 := NewTime(0, 0, 0, 0)

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

func TestTime_NewTimeFromString(t *testing.T) {
	noMs := "23:07:10"
	expNoMs := NewTime(23, 7, 10, 0)
	timeNoMs, err := NewTimeFromString(noMs)
	if err != nil {
		t.Error(err)
	}
	if !expNoMs.equals(&timeNoMs) {
		t.Error(timeNoMs)
	}

	ms := "23:07:10.957"
	expMs := NewTime(23, 7, 10, 957)
	timeMs, err := NewTimeFromString(ms)
	if err != nil {
		t.Error(err)
	}
	if !expMs.equals(&timeMs) {
		t.Error(timeMs)
	}
}

func TestTime_ToZinc(t *testing.T) {
	timeNoMs := NewTime(23, 7, 10, 0)
	timeNoMsZinc := timeNoMs.ToZinc()
	if timeNoMsZinc != "23:07:10" {
		t.Error(timeNoMsZinc)
	}

	timeMs := NewTime(23, 7, 10, 957)
	timeMsZinc := timeMs.ToZinc()
	if timeMsZinc != "23:07:10.957" {
		t.Error(timeMsZinc)
	}

	timeOnesMs := NewTime(23, 7, 10, 2)
	timeOnesMsZinc := timeOnesMs.ToZinc()
	if timeOnesMsZinc != "23:07:10.002" {
		t.Error(timeOnesMsZinc)
	}

	timeTensMs := NewTime(23, 7, 10, 56)
	timeTensMsZinc := timeTensMs.ToZinc()
	if timeTensMsZinc != "23:07:10.056" {
		t.Error(timeTensMsZinc)
	}
}
