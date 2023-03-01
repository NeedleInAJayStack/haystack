package haystack

import (
	"encoding/json"
	"testing"
)

func TestDate_equals(t *testing.T) {
	date1 := NewDate(2020, 8, 17)
	date2 := NewDate(2020, 8, 17)
	date3 := NewDate(0, 0, 0)

	if !date1.equals(date1) {
		t.Error("The same object doesn't equal itself")
	}
	if !date1.equals(date2) {
		t.Error("Equivalent objects doesn't equal itself")
	}
	if !date2.equals(date1) {
		t.Error("Ordering matters")
	}
	if date1.equals(date3) {
		t.Error("Non-equivalent objects are equal")
	}
}

func TestDate_NewDateFromIso(t *testing.T) {
	dateStr := "2020-08-17"
	exp := NewDate(2020, 8, 17)
	date, err := NewDateFromIso(dateStr)
	if err != nil {
		t.Error(err)
	}
	if !exp.equals(date) {
		t.Error(date)
	}
}

func TestDate_ToZinc(t *testing.T) {
	valTest_ToZinc(NewDate(2020, 8, 17), "2020-08-17", t)
}

func TestDate_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewDate(2020, 8, 17), "\"d:2020-08-17\"", t)
}

func TestDate_UnmarshalJSON(t *testing.T) {
	var date Date
	json.Unmarshal([]byte("\"d:2020-08-17\""), &date)
	valTest_ToZinc(date, "2020-08-17", t)
}

func TestDate_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewDate(2020, 8, 17), "{\"_kind\":\"date\",\"val\":\"2020-08-17\"}", t)
}
