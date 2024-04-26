package haystack

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime_NewDateTimeRaw(t *testing.T) {
	utc, err := NewDateTimeRaw(2020, 8, 17, 23, 07, 10, 0, "UTC")
	assert.Nil(t, err)
	assert.Equal(t, 2020, utc.time.Year())
	assert.Equal(t, time.August, utc.time.Month())
	assert.Equal(t, 17, utc.time.Day())
	assert.Equal(t, 7, utc.time.Minute())
	assert.Equal(t, 10, utc.time.Second())
	assert.Equal(t, 0, utc.time.Nanosecond())
	assert.Equal(t, "UTC", utc.time.Location().String())
}

func TestDateTime_NewDateTimeFromString(t *testing.T) {
	utc, err := NewDateTimeFromString("2020-08-17T23:07:10Z UTC")
	assert.Nil(t, err)
	assert.Equal(t, 2020, utc.time.Year())
	assert.Equal(t, time.August, utc.time.Month())
	assert.Equal(t, 17, utc.time.Day())
	assert.Equal(t, 7, utc.time.Minute())
	assert.Equal(t, 10, utc.time.Second())
	assert.Equal(t, 0, utc.time.Nanosecond())
	assert.Equal(t, "UTC", utc.Tz())

	la, _ := NewDateTimeFromString("2020-08-17T23:07:10-07:00 Los_Angeles")
	assert.Equal(t, "Los_Angeles", la.Tz())

	taipei, _ := NewDateTimeFromString("2020-08-17T23:07:10+08:00 Taipei")
	assert.Equal(t, "Taipei", taipei.Tz())
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
