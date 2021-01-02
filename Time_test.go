package haystack

import "testing"

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
	noMs := "23:07:10"
	expNoMs := NewTime(23, 7, 10, 0)
	timeNoMs, err := NewTimeFromIso(noMs)
	if err != nil {
		t.Error(err)
	}
	if !expNoMs.equals(timeNoMs) {
		t.Error(timeNoMs)
	}

	ms := "23:07:10.957"
	expMs := NewTime(23, 7, 10, 957)
	timeMs, err := NewTimeFromIso(ms)
	if err != nil {
		t.Error(err)
	}
	if !expMs.equals(timeMs) {
		t.Error(timeMs)
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
	noMs := NewTime(0, 0, 0, 0)
	valTest_UnmarshalJSON("\"h:23:07:10\"", noMs, "23:07:10", t)
	oneMs := NewTime(0, 0, 0, 0)
	valTest_UnmarshalJSON("\"h:23:07:10.002\"", oneMs, "23:07:10.002", t)
	tenMs := NewTime(0, 0, 0, 0)
	valTest_UnmarshalJSON("\"h:23:07:10.056\"", tenMs, "23:07:10.056", t)
	hundredMs := NewTime(0, 0, 0, 0)
	valTest_UnmarshalJSON("\"h:23:07:10.957\"", hundredMs, "23:07:10.957", t)
}

func TestTime_MarshalHayson(t *testing.T) {
	valTest_MarshalHayson(NewTime(23, 7, 10, 0), "{\"_kind\":\"time\",\"val\":\"23:07:10\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 2), "{\"_kind\":\"time\",\"val\":\"23:07:10.002\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 56), "{\"_kind\":\"time\",\"val\":\"23:07:10.056\"}", t)
	valTest_MarshalHayson(NewTime(23, 7, 10, 957), "{\"_kind\":\"time\",\"val\":\"23:07:10.957\"}", t)
}
