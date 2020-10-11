package haystack

// TODO: Determine why we get "400 Bad Request" on BuildingFit SkySpark instance
// 		- My best guess is that it's an NGINX issue
// TODO: Understand/document "x509: certificate signed by unknown authority" - Only happens on linux
//		- https://bugs.launchpad.net/snappy/+bug/1620755/
//		- https://github.com/golang/go/issues/12139

// func TestClient(t *testing.T) {

// 	// Test with local SkySpark instance
// 	haystackClient := NewClient(
// 		"http://localhost:8080/api/demo",
// 		"test",
// 		"test",
// 	)

// 	openErr := haystackClient.Open()
// 	if openErr != nil {
// 		t.Error(openErr)
// 	}

// 	_, callErr := haystackClient.Call("about", EmptyGrid())
// 	if callErr != nil {
// 		t.Error(callErr)
// 	}

// 	_, aboutErr := haystackClient.About()
// 	if aboutErr != nil {
// 		t.Error(aboutErr)
// 	}

// 	_, formatsErr := haystackClient.Formats()
// 	if formatsErr != nil {
// 		t.Error(formatsErr)
// 	}

// 	_, opsErr := haystackClient.Ops()
// 	if opsErr != nil {
// 		t.Error(opsErr)
// 	}

// 	_, readErr := haystackClient.Read("site")
// 	if readErr != nil {
// 		t.Error(readErr)
// 	}

// 	readLimit, readLimitErr := haystackClient.ReadLimit("point", 1)
// 	if readLimitErr != nil {
// 		t.Error(readLimitErr)
// 	}

// 	pointRef := readLimit.RowAt(0).Get("id").(Ref)

// 	_, readByIdsErr := haystackClient.ReadByIds([]Ref{pointRef})
// 	if readByIdsErr != nil {
// 		t.Error(readByIdsErr)
// 	}

// 	_, hisReadErr := haystackClient.HisRead(pointRef, "yesterday")
// 	if hisReadErr != nil {
// 		t.Error(hisReadErr)
// 	}

// 	fromDate := NewDate(2020, 10, 4)
// 	toDate := NewDate(2020, 10, 5)
// 	_, hisReadAbsDateErr := haystackClient.HisReadAbsDate(pointRef, fromDate, toDate)
// 	if hisReadAbsDateErr != nil {
// 		t.Error(hisReadAbsDateErr)
// 	}

// 	fromTs, _ := NewDateTimeFromString("2020-10-04T00:00:00-07:00 Los_Angeles")
// 	toTs, _ := NewDateTimeFromString("2020-10-05T00:00:00-07:00 Los_Angeles")
// 	_, hisReadAbsDateTimeErr := haystackClient.HisReadAbsDateTime(pointRef, fromTs, toTs)
// 	if hisReadAbsDateTimeErr != nil {
// 		t.Error(hisReadAbsDateTimeErr)
// 	}

// 	_, evalErr := haystackClient.Eval("read(point)")
// 	if evalErr != nil {
// 		t.Error(evalErr)
// 	}
// }
