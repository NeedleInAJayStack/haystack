package haystack

import (
	"testing"
	"time"
)

func TestDateTime_toZinc(t *testing.T) {
	utc := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC))
	utcZinc := utc.toZinc()
	if utcZinc != "2020-08-17T23:07:10Z UTC" {
		t.Error(utcZinc)
	}

	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	la := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc))
	laZinc := la.toZinc()
	if laZinc != "2020-08-17T23:07:10-07:00 Los_Angeles" {
		t.Error(laZinc)
	}

	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")
	taipei := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc))
	taipeiZinc := taipei.toZinc()
	if taipeiZinc != "2020-08-17T23:07:10+08:00 Taipei" {
		t.Error(taipeiZinc)
	}
}
