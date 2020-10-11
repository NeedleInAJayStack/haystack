package haystack

import (
	"testing"
	"time"
)

func TestDateTime_NewDateTime(t *testing.T) {
	utc := NewDateTime(2020, 8, 17, 23, 07, 10, 0, 0, "UTC")
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
		t.Error(utc.time.min)
	}
	if utc.time.sec != 10 {
		t.Error(utc.time.sec)
	}
	if utc.time.ms != 0 {
		t.Error(utc.time.ms)
	}
	if utc.tzOffset != 0 {
		t.Error(utc.tzOffset)
	}
	if utc.tz != "UTC" {
		t.Error(utc.tz)
	}
}

func TestDateTime_NewDateTimeFromString(t *testing.T) {
	utc, _ := NewDateTimeFromString("2020-08-17T23:07:10Z UTC")
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
		t.Error(utc.time.min)
	}
	if utc.time.sec != 10 {
		t.Error(utc.time.sec)
	}
	if utc.time.ms != 0 {
		t.Error(utc.time.ms)
	}
	if utc.tzOffset != 0 {
		t.Error(utc.tzOffset)
	}
	if utc.tz != "UTC" {
		t.Error(utc.tz)
	}

	la, _ := NewDateTimeFromString("2020-08-17T23:07:10-07:00 Los_Angeles")
	if la.tzOffset != -25200 {
		t.Error(la.tzOffset)
	}
	if la.tz != "Los_Angeles" {
		t.Error(la.tz)
	}

	taipei, _ := NewDateTimeFromString("2020-08-17T23:07:10+08:00 Taipei")
	if taipei.tzOffset != 28800 {
		t.Error(la.tzOffset)
	}
	if taipei.tz != "Taipei" {
		t.Error(la.tz)
	}
}

func TestDateTime_ToZinc(t *testing.T) {
	utc := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC))
	valTest_ToZinc(utc, "2020-08-17T23:07:10Z UTC", t)

	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	la := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc))
	valTest_ToZinc(la, "2020-08-17T23:07:10-07:00 Los_Angeles", t)

	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")
	taipei := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc))
	valTest_ToZinc(taipei, "2020-08-17T23:07:10+08:00 Taipei", t)
}

func TestDateTime_MarshalJSON(t *testing.T) {
	utc := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC))
	valTest_MarshalJSON(utc, "\"t:2020-08-17T23:07:10Z UTC\"", t)

	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	la := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc))
	valTest_MarshalJSON(la, "\"t:2020-08-17T23:07:10-07:00 Los_Angeles\"", t)

	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")
	taipei := dateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc))
	valTest_MarshalJSON(taipei, "\"t:2020-08-17T23:07:10+08:00 Taipei\"", t)
}
