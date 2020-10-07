package haystack

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {

	// Test with local SkySpark instance
	uri := "http://localhost:8080/api/demo"
	user := "test"
	pass := "test"

	haystackClient := NewClient(uri, user, pass)
	openErr := haystackClient.Open()
	if openErr != nil {
		t.Error(openErr)
	}

	call, callErr := haystackClient.Call("about", haystack.EmptyGrid())
	if callErr != nil {
		fmt.Println(call.ToZinc())
		t.Error(callErr)
	}

	about, aboutErr := haystackClient.About()
	if aboutErr != nil {
		fmt.Println(about.ToZinc())
		t.Error(aboutErr)
	}

	formats, formatsErr := haystackClient.Formats()
	if formatsErr != nil {
		fmt.Println(formats.ToZinc())
		t.Error(formatsErr)
	}

	ops, opsErr := haystackClient.Ops()
	if opsErr != nil {
		fmt.Println(ops.ToZinc())
		t.Error(opsErr)
	}

	read, readErr := haystackClient.Read("site")
	if readErr != nil {
		fmt.Println(read.ToZinc())
		t.Error(readErr)
	}

	readLimit, readLimitErr := haystackClient.ReadLimit("point", 1)
	if readLimitErr != nil {
		fmt.Println(readLimit.ToZinc())
		t.Error(readLimitErr)
	}

	hisRef := readLimit.RowAt(0).Get("id").(haystack.Ref)

	hisRead, hisReadErr := haystackClient.HisRead(hisRef, "yesterday")
	if hisReadErr != nil {
		fmt.Println(hisRead.ToZinc())
		t.Error(hisReadErr)
	}

	fromDate := haystack.NewDate(2020, 10, 4)
	toDate := haystack.NewDate(2020, 10, 5)
	hisReadAbsDate, hisReadAbsDateErr := haystackClient.HisReadAbsDate(hisRef, fromDate, toDate)
	if hisReadAbsDateErr != nil {
		fmt.Println(hisReadAbsDate.ToZinc())
		t.Error(hisReadAbsDateErr)
	}

	fromTs, _ := haystack.NewDateTimeFromString("2020-10-04T00:00:00-07:00 Los_Angeles")
	toTs, _ := haystack.NewDateTimeFromString("2020-10-05T00:00:00-07:00 Los_Angeles")
	hisReadAbsDateTime, hisReadAbsDateTimeErr := haystackClient.HisReadAbsDateTime(hisRef, fromTs, toTs)
	if hisReadAbsDateTimeErr != nil {
		fmt.Println(hisReadAbsDateTime.ToZinc())
		t.Error(hisReadAbsDateTimeErr)
	}

	eval, evalErr := haystackClient.Eval("read(point)")
	if evalErr != nil {
		fmt.Println(eval.ToZinc())
		t.Error(evalErr)
	}
}
