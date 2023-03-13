package haystack

import (
	"encoding/json"
	"testing"
)

func TestTime_equals(t *testing.T) {
	time1 := NewTime(23, 7, 10, 957)
	time2 := NewTime(23, 7, 10, 957)
	time3 := NewTime(0, 0, 0, 0)

	if !time1.equals(time1) {
		t.Error("The same object doesn't equal itself")
	}
	if !time1.equals(time2) {
		t.Error("Equivalent objects doesn't equal itself")
	}
	if !time2.equals(time1) {
		t.Error("Ordering matters")
	}
	if time1.equals(time3) {
		t.Error("Non-equivalent objects are equal")
	}
}

func TestTime_NewTimeFromIso(t *testing.T) {
	iso := "23:07:10"
	expected := NewTime(23, 7, 10, 0)
	actual, err := NewTimeFromIso(iso)
	if err != nil {
		t.Error(err)
	}
	if !expected.equals(actual) {
		t.Error(actual)
	}
}

func TestTime_NewTimeFromIso_ms(t *testing.T) {
	iso := "23:07:10.957"
	expected := NewTime(23, 7, 10, 957)
	actual, err := NewTimeFromIso(iso)
	if err != nil {
		t.Error(err)
	}
	if !expected.equals(actual) {
		t.Error(actual)
	}
}

func TestTime_NewTimeFromIso_manyMs(t *testing.T) {
	iso := "23:07:10.957654321"
	expected := NewTime(23, 7, 10, 957)
	actual, err := NewTimeFromIso(iso)
	if err != nil {
		t.Error(err)
	}
	if !expected.equals(actual) {
		t.Error(actual)
	}
}

func TestTime_ToZinc(t *testing.T) {
	valTest_ToZinc(NewTime(23, 7, 10, 0), "23:07:10", t)
	valTest_ToZinc(NewTime(23, 7, 10, 2), "23:07:10.002", t)
	valTest_ToZinc(NewTime(23, 7, 10, 56), "23:07:10.056", t)
	valTest_ToZinc(NewTime(23, 7, 10, 957), "23:07:10.957", t)
}

func TestTime_MarshalJSON(t *testing.T) {
	valTest_MarshalJSON(NewTime(23, 7, 10, 0), "\"h:23:07:10\"", t)
	valTest_MarshalJSON(NewTime(23, 7, 10, 2), "\"h:23:07:10.002\"", t)
	valTest_MarshalJSON(NewTime(23, 7, 10, 56), "\"h:23:07:10.056\"", t)
	valTest_MarshalJSON(NewTime(23, 7, 10, 957), "\"h:23:07:10.957\"", t)
}

func TestTime_UnmarshalJSON(t *testing.T) {
	var noMs Time
	json.Unmarshal([]byte("\"h:23:07:10\""), &noMs)
	valTest_ToZinc(noMs, "23:07:10", t)

	var oneMs Time
	json.Unmarshal([]byte("\"h:23:07:10.002\""), &oneMs)
	valTest_ToZinc(oneMs, "23:07:10.002", t)

	var tenMs Time
	json.Unmarshal([]byte("\"h:23:07:10.056\""), &tenMs)
	valTest_ToZinc(tenMs, "23:07:10.056", t)

	var hundredMs Time
	json.Unmarshal([]byte("\"h:23:07:10.957\""), &hundredMs)
	valTest_ToZinc(hundredMs, "23:07:10.957", t)
}

func TestTime_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewTime(23, 7, 10, 0), "{\"_kind\":\"time\",\"val\":\"23:07:10\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 2), "{\"_kind\":\"time\",\"val\":\"23:07:10.002\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 56), "{\"_kind\":\"time\",\"val\":\"23:07:10.056\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 957), "{\"_kind\":\"time\",\"val\":\"23:07:10.957\"}", t)
}
