package haystack

import "testing"

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
