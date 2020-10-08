package haystack

import (
	"testing"
)

func TestClient(t *testing.T) {

	// Test with local SkySpark instance
	haystackClient := NewClient(
		"http://localhost:8080/api/demo",
		"test",
		"test",
	)
	openErr := haystackClient.Open()
	if openErr != nil {
		t.Error(openErr)
	}

	_, callErr := haystackClient.Call("about", EmptyGrid())
	if callErr != nil {
		t.Error(callErr)
	}

	_, aboutErr := haystackClient.About()
	if aboutErr != nil {
		t.Error(aboutErr)
	}

	_, formatsErr := haystackClient.Formats()
	if formatsErr != nil {
		t.Error(formatsErr)
	}

	_, opsErr := haystackClient.Ops()
	if opsErr != nil {
		t.Error(opsErr)
	}

	_, readErr := haystackClient.Read("site")
	if readErr != nil {
		t.Error(readErr)
	}

	readLimit, readLimitErr := haystackClient.ReadLimit("point", 1)
	if readLimitErr != nil {
		t.Error(readLimitErr)
	}

	hisRef := readLimit.RowAt(0).Get("id").(Ref)

	_, hisReadErr := haystackClient.HisRead(hisRef, "yesterday")
	if hisReadErr != nil {
		t.Error(hisReadErr)
	}

	fromDate := NewDate(2020, 10, 4)
	toDate := NewDate(2020, 10, 5)
	_, hisReadAbsDateErr := haystackClient.HisReadAbsDate(hisRef, fromDate, toDate)
	if hisReadAbsDateErr != nil {
		t.Error(hisReadAbsDateErr)
	}

	fromTs, _ := NewDateTimeFromString("2020-10-04T00:00:00-07:00 Los_Angeles")
	toTs, _ := NewDateTimeFromString("2020-10-05T00:00:00-07:00 Los_Angeles")
	_, hisReadAbsDateTimeErr := haystackClient.HisReadAbsDateTime(hisRef, fromTs, toTs)
	if hisReadAbsDateTimeErr != nil {
		t.Error(hisReadAbsDateTimeErr)
	}

	_, evalErr := haystackClient.Eval("read(point)")
	if evalErr != nil {
		t.Error(evalErr)
	}
}
