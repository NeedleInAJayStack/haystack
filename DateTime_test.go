package haystack

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_NewDateTimeRaw(t *testing.T) {
	utc := NewDateTimeRaw(2020, 8, 17, 23, 07, 10, 0, 0, "UTC")
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
	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	assert.Equal(t, NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC)).ToZinc(), "2020-08-17T23:07:10Z UTC")
	assert.Equal(t, NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc)).ToZinc(), "2020-08-17T23:07:10-07:00 Los_Angeles")
	assert.Equal(t, NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc)).ToZinc(), "2020-08-17T23:07:10+08:00 Taipei")
}

func TestDateTime_ToAxon(t *testing.T) {
	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	valTest_ToAxon(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC)), "dateTime(2020-08-17,23:07:10,\"UTC\")", t)
	valTest_ToAxon(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc)), "dateTime(2020-08-17,23:07:10,\"Los_Angeles\")", t)
	valTest_ToAxon(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc)), "dateTime(2020-08-17,23:07:10,\"Taipei\")", t)
}

func TestDateTime_MarshalJSON(t *testing.T) {
	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	valTest_MarshalJSON(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC)), "\"t:2020-08-17T23:07:10Z UTC\"", t)
	valTest_MarshalJSON(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc)), "\"t:2020-08-17T23:07:10-07:00 Los_Angeles\"", t)
	valTest_MarshalJSON(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc)), "\"t:2020-08-17T23:07:10+08:00 Taipei\"", t)
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	var utc DateTime
	json.Unmarshal([]byte("\"t:2020-08-17T23:07:10Z UTC\""), &utc)
	assert.Equal(t, utc.ToZinc(), "2020-08-17T23:07:10Z UTC")

	var la DateTime
	json.Unmarshal([]byte("\"t:2020-08-17T23:07:10-07:00 Los_Angeles\""), &la)
	assert.Equal(t, la.ToZinc(), "2020-08-17T23:07:10-07:00 Los_Angeles")

	var taipei DateTime
	json.Unmarshal([]byte("\"t:2020-08-17T23:07:10+08:00 Taipei\""), &taipei)
	assert.Equal(t, taipei.ToZinc(), "2020-08-17T23:07:10+08:00 Taipei")
}

func TestDateTime_MarshalHayson(t *testing.T) {
	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	valTest_MarshalHayson(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC)), "{\"_kind\":\"dateTime\",\"val\":\"2020-08-17T23:07:10Z\",\"tz\":\"UTC\"}", t)
	valTest_MarshalHayson(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc)), "{\"_kind\":\"dateTime\",\"val\":\"2020-08-17T23:07:10-07:00\",\"tz\":\"Los_Angeles\"}", t)
	valTest_MarshalHayson(NewDateTimeFromGo(time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc)), "{\"_kind\":\"dateTime\",\"val\":\"2020-08-17T23:07:10+08:00\",\"tz\":\"Taipei\"}", t)
}

func TestDateTime_ToGo(t *testing.T) {
	losAngelesLoc, _ := time.LoadLocation("America/Los_Angeles")
	taipeiLoc, _ := time.LoadLocation("Asia/Taipei")

	date1 := time.Date(2020, time.August, 17, 23, 7, 10, 0, time.UTC)
	if date1.Unix() != NewDateTimeFromGo(date1).ToGo().Unix() {
		t.Error(date1, "!=", NewDateTimeFromGo(date1).ToGo())
	}

	date2 := time.Date(2020, time.August, 17, 23, 7, 10, 0, losAngelesLoc)
	if date2.Unix() != NewDateTimeFromGo(date2).ToGo().Unix() {
		t.Error(date2, "!=", NewDateTimeFromGo(date2).ToGo())
	}

	date3 := time.Date(2020, time.August, 17, 23, 7, 10, 0, taipeiLoc)
	if date3.Unix() != NewDateTimeFromGo(date3).ToGo().Unix() {
		t.Error(date3, "!=", NewDateTimeFromGo(date3).ToGo())
	}

	date4 := time.Date(2020, time.August, 17, 23, 7, 10, 0, time.Local)
	if date4.Unix() != NewDateTimeFromGo(date4).ToGo().Unix() {
		t.Error(date4, "!=", NewDateTimeFromGo(date4).ToGo())
	}
}

func valTest_ToAxon(val DateTime, expected string, t *testing.T) {
	actual := val.ToAxon()
	if actual != expected {
		t.Error(actual + " != " + expected)
	}
}
