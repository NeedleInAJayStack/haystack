package haystack

import (
	"testing"
	"time"
)

func TestDateTime_dateTimefromStr(t *testing.T) {
	utc, _ := dateTimeFromStr("2020-08-17T23:07:10Z UTC")
	if utc.date.year != 2020 {
		t.Error(utc.date.year)
	}
	if utc.date.month != 8 {
		t.Error(utc.date.month)
	}
	if utc.date.day != 17 {
		t.Error(utc.date.day)
	}
	if utc.time.hour != 23 {
		t.Error(utc.time.hour)
	}
	if utc.time.min != 7 {
		t.Error(utc.time.hour)
	}
	if utc.tzOffset != 0 {
		t.Error(utc.tzOffset)
	}
	if utc.tz != "UTC" {
		t.Error(utc.tz)
	}

	la, _ := dateTimeFromStr("2020-08-17T23:07:10-07:00 Los_Angeles")
	if la.tzOffset != -25200 {
		t.Error(la.tzOffset)
	}
	if la.tz != "Los_Angeles" {
		t.Error(la.tz)
	}

	taipei, _ := dateTimeFromStr("2020-08-17T23:07:10+08:00 Taipei")
	if taipei.tzOffset != 28800 {
		t.Error(la.tzOffset)
	}
	if taipei.tz != "Taipei" {
		t.Error(la.tz)
	}
}

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
