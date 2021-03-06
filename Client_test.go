package haystack

import (
	"errors"
	"testing"
)

func TestClient_Open(t *testing.T) {
	client := testClient()
	openErr := client.Open()
	if openErr != nil {
		t.Error(openErr)
	}
}

func TestClient_Call(t *testing.T) {
	client := testClient()
	actual, err := client.Call("about", EmptyGrid())
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_about, t)
}

func TestClient_About(t *testing.T) {
	client := testClient()
	actual, aboutErr := client.About()
	if aboutErr != nil {
		t.Error(aboutErr)
	}
	// About returns a Dict so get the first value of expected manually
	expectedZinc := clientHTTPMock_about
	var reader ZincReader
	reader.InitString(expectedZinc)
	expectedGrid := reader.ReadVal()
	expected := expectedGrid.(*Grid).RowAt(0).ToDict()

	valTest_Equal_Grid(actual, expected, t)
}

func TestClient_Formats(t *testing.T) {
	client := testClient()
	actual, err := client.Formats()
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_formats, t)
}

func TestClient_Ops(t *testing.T) {
	client := testClient()
	actual, err := client.Ops()
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_ops, t)
}

func TestClient_Read(t *testing.T) {
	client := testClient()
	actual, err := client.Read("site")
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_readSites, t)
}

func TestClient_ReadLimit(t *testing.T) {
	client := testClient()
	actual, err := client.ReadLimit("point", 1)
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_readPoint, t)
}

func TestClient_ReadByIds(t *testing.T) {
	client := testClient()
	points, pointsErr := client.ReadLimit("point", 1)
	if pointsErr != nil {
		t.Error(pointsErr)
	} else {
		pointRef := points.RowAt(0).Get("id").(*Ref)

		actual, err := client.ReadByIds([]*Ref{pointRef})
		if err != nil {
			t.Error(err)
		}
		testClient_ValZinc(actual, clientHTTPMock_readPoint, t)
	}
}

func TestClient_HisRead(t *testing.T) {
	client := testClient()
	points, pointsErr := client.ReadLimit("point", 1)
	if pointsErr != nil {
		t.Error(pointsErr)
	} else {
		pointRef := points.RowAt(0).Get("id").(*Ref)

		actual, err := client.HisRead(pointRef, "yesterday")
		if err != nil {
			t.Error(err)
		}
		testClient_ValZinc(actual, clientHTTPMock_hisRead20210103, t)
	}
}

func TestClient_HisReadAbsDate(t *testing.T) {
	client := testClient()
	points, pointsErr := client.ReadLimit("point", 1)
	if pointsErr != nil {
		t.Error(pointsErr)
	} else {
		pointRef := points.RowAt(0).Get("id").(*Ref)

		fromDate := NewDate(2020, 10, 4)
		toDate := NewDate(2020, 10, 5)
		actual, err := client.HisReadAbsDate(pointRef, fromDate, toDate)
		if err != nil {
			t.Error(err)
		}
		testClient_ValZinc(actual, clientHTTPMock_hisRead20201004to6, t)
	}
}

func TestClient_HisReadAbsDateTime(t *testing.T) {
	client := testClient()
	points, pointsErr := client.ReadLimit("point", 1)
	if pointsErr != nil {
		t.Error(pointsErr)
	} else {
		pointRef := points.RowAt(0).Get("id").(*Ref)

		fromTs, _ := NewDateTimeFromString("2020-10-04T00:00:00-07:00 Los_Angeles")
		toTs, _ := NewDateTimeFromString("2020-10-05T00:00:00-07:00 Los_Angeles")
		actual, err := client.HisReadAbsDateTime(pointRef, fromTs, toTs)
		if err != nil {
			t.Error(err)
		}
		testClient_ValZinc(actual, clientHTTPMock_hisReadDateTimes, t)
	}
}

func TestClient_Eval(t *testing.T) {
	client := testClient()
	actual, err := client.Eval("read(point)")
	if err != nil {
		t.Error(err)
	}
	testClient_ValZinc(actual, clientHTTPMock_readPoint, t)
}

func testClient_ValZinc(actual Val, expectedZinc string, t *testing.T) {
	var reader ZincReader
	reader.InitString(expectedZinc)
	expected := reader.ReadVal()
	valTest_Equal_Grid(actual, expected, t)
}

func testClient() *Client {
	return &Client{
		clientHTTP: &clientHTTPMock{},
		uri:        "http://localhost:8080/api/demo",
		username:   "test",
		password:   "test",
	}
}

// clientHTTPMock allows us to remove the HTTP dependency within tests
type clientHTTPMock struct {
}

func (clientHTTPMock *clientHTTPMock) open(uri string, username string, password string) (string, error) {
	// For now, just say we did it
	return "test", nil
}

func (clientHTTPMock *clientHTTPMock) postString(uri string, auth string, op string, reqBody string) (string, error) {
	// These are taken from a SkySpark 3.0.26 demo project on 2021-01-03
	if op == "about" {
		// Can't use string literal because of Uri backticks
		return clientHTTPMock_about, nil
	} else if op == "formats" {
		return clientHTTPMock_formats, nil
	} else if op == "ops" {
		return clientHTTPMock_ops, nil
	} else if op == "read" {
		if reqBody == "ver:\"3.0\"\nfilter, limit\n\"site\", N" { // readAll sites
			return clientHTTPMock_readSites, nil
		} else if reqBody == "ver:\"3.0\"\nfilter, limit\n\"point\", 1" { // readLimit point
			return clientHTTPMock_readPoint, nil
		} else if reqBody == "ver:\"3.0\"\nid\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" { // readById
			return clientHTTPMock_readPoint, nil
		}
		return emptyRes, errors.New("'read' argument not supported by mock class")
	} else if op == "hisRead" {
		if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"yesterday\"" { // hisRead relative
			return clientHTTPMock_hisRead20210103, nil
		} else if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"2020-10-04,2020-10-05\"" { // hisRead absolute dates
			return clientHTTPMock_hisRead20201004to6, nil
		} else if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"2020-10-04T00:00:00-07:00 Los_Angeles,2020-10-05T00:00:00-07:00 Los_Angeles\"" { // hisRead absolute datetimes
			return clientHTTPMock_hisReadDateTimes, nil
		}
		return emptyRes, errors.New("'hisRead' argument not supported by mock class")
	} else if op == "eval" {
		if reqBody == "ver:\"3.0\"\nexpr\n\"read(point)\"" { // eval read(point)
			return clientHTTPMock_readPoint, nil
		}
		return emptyRes, errors.New("'eval' argument not supported by mock class")
	}
	return emptyRes, errors.New("haystack op not supported by mock class: " + op)
}

const (
	clientHTTPMock_about string = "ver:\"3.0\"\n" + // Can't use string literal because of Uri backticks
		"haystackVersion,projName,serverName,serverBootTime,serverTime,productName,productUri,productVersion,moduleName,moduleVersion,tz,whoami,hostDis,hostModel,hostId\n" +
		"\"3.0\",\"demo\",\"JaysDesktop\",2021-01-03T00:21:01.588-07:00 Denver,2021-01-03T00:21:43.799-07:00 Denver,\"SkySpark\",`http://skyfoundry.com/skyspark`,\"3.0.26\",\"skyarcd\",\"3.0.26\",\"Denver\",\"test\",\"Linux amd64 5.4.0-58-generic\",\"Linux amd64 5.4.0-58-generic\",NA\n"
	clientHTTPMock_formats string = `ver:"3.0"
		name,dis,mime,receive,send
		"zinc","Zinc","text/plain",M,M
		"csv","CSV","text/csv",M,M
		"excel","Excel","application/vnd.ms-excel",,M
		"json","JSON","application/json",M,M
		"jsonld","JSON-LD","application/ld+json",,M
		"pdf","PDF","application/pdf",,M
		"svg","SVG","image/svg+xml",,M
		"trio","Trio","text/trio",M,M
		"turtle","Turtle","text/turtle",,M
		"xml","XML","text/xml",,M
		"zinc","Zinc","text/zinc",M,M
		`
	clientHTTPMock_ops string = `ver:"3.0"
		name,summary
		"about","Summary info for server"
		"commit","Commit diffs to proj database"
		"eval","Evaluate an Axon expression"
		"formats","Data formats supported by server"
		"hisRead","Read time series data to historian"
		"hisWrite","Write time series data to historian"
		"invokeAction","Watch subscription"
		"nav","Learn navigation"
		"ops","Operations supported by server"
		"pointWrite","Read/write point write array"
		"read","Read records by id or filter"
		"watchPoll","Watch poll cov or refresh"
		"watchSub","Watch subscription"
		"watchUnsub","Watch unsubscription"
		`
	clientHTTPMock_readSites string = `ver:"3.0"
		id,area,dis,geoAddr,geoCity,geoCoord,geoCountry,geoPostalCode,geoState,geoStreet,hq,metro,occupiedEnd,occupiedStart,primaryFunction,regionRef,site,store,storeNum,tz,weatherStationRef,yearBuilt,mod
		@p:demo:r:2725da26-ac563571 "Headquarters",140797ft??,"Headquarters","600 W Main St, Richmond, VA","Richmond",C(37.545826,-77.449188),"US","23220","VA","600 W Main St",M,"Richmond",18:00:00,08:00:00,"Office",@p:demo:r:2725da26-f3e488bc "Richmond",M,,,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1999,2020-10-23T18:15:02.701Z
		@p:demo:r:2725da26-3ca6125c "Gaithersburg",8013ft??,"Gaithersburg","18212 Montgomery Village Ave, Gaithersburg, MD","Gaithersburg",C(39.154824,-77.209002),"US","20879","MD","18212 Montgomery Village Ave",,"Washington DC",21:00:00,09:00:00,"Retail Store",@p:demo:r:2725da26-e77a16f1 "Washington DC",M,M,4,"New_York",@p:demo:r:2725da26-9bb170b8 "Washington, DC",2001,2020-10-23T18:15:02.797Z
		@p:demo:r:2725da26-d280b1b5 "Short Pump",17122ft??,"Short Pump","11282 W Broad St, Richmond, VA","Glen Allen",C(37.650338,-77.606105),"US","23060","VA","11282 W Broad St",,"Richmond",21:00:00,10:00:00,"Retail Store",@p:demo:r:2725da26-f3e488bc "Richmond",M,M,3,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1999,2020-10-23T18:15:02.763Z
		@p:demo:r:2725da26-505b4ae8 "Carytown",3149ft??,"Carytown","3504 W Cary St, Richmond, VA","Richmond",C(37.555385,-77.486903),"US","23221","VA","3504 W Cary St",,"Richmond",20:00:00,10:00:00,"Retail Store",@p:demo:r:2725da26-f3e488bc "Richmond",M,M,1,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1996,2020-10-23T18:15:02.742Z
		`
	clientHTTPMock_readPoint string = `ver:"3.0"
		id,navName,disMacro,point,his,siteRef,equipRef,curVal,curStatus,hisEnd,hisSize,hisStart,kind,tz,cmd,elecRef,cur,regionRef,fan,discharge,air,hisMode,enum,mod
		@p:demo:r:2725da26-1dda68ee "Gaithersburg RTU-1 Fan","Fan","\\$equipRef \\$navName",M,M,@p:demo:r:2725da26-3ca6125c "Gaithersburg",@p:demo:r:2725da26-1d885a68 "Gaithersburg RTU-1",F,"ok",2021-01-21T22:15:00-05:00 New_York,2625,2019-01-01T00:00:00-05:00 New_York,"Bool","New_York",M,@p:demo:r:2725da26-8ddc7cf5 "Gaithersburg ElecMeter-Hvac",M,@p:demo:r:2725da26-e77a16f1 "Washington DC",M,M,M,"cov","off,on",2020-10-23T18:15:02.83Z
		`
	clientHTTPMock_hisRead20210103 string = `ver:"3.0" id:@p:demo:r:2725da26-1dda68ee "Gaithersburg RTU-1 Fan" hisStart:2021-01-02T00:00:00-05:00 New_York hisEnd:2021-01-03T00:00:00-05:00 New_York
		ts,val
		2021-01-02T00:00:00-05:00 New_York,F
		2021-01-02T09:00:00-05:00 New_York,T
		2021-01-02T21:15:00-05:00 New_York,F
		`
	clientHTTPMock_hisRead20201004to6 string = `ver:"3.0" id:@p:demo:r:2725da26-1dda68ee "Gaithersburg RTU-1 Fan" hisStart:2020-10-04T00:00:00-04:00 New_York hisEnd:2020-10-06T00:00:00-04:00 New_York
		ts,val
		2020-10-04T00:00:00-04:00 New_York,F
		2020-10-04T09:00:00-04:00 New_York,T
		2020-10-04T21:15:00-04:00 New_York,F
		2020-10-05T00:00:00-04:00 New_York,F
		2020-10-05T09:00:00-04:00 New_York,T
		2020-10-05T21:15:00-04:00 New_York,F
		`
	clientHTTPMock_hisReadDateTimes string = `ver:"3.0" id:@p:demo:r:2725da26-1dda68ee "Gaithersburg RTU-1 Fan" hisStart:2020-10-04T03:00:00-04:00 New_York hisEnd:2020-10-05T03:00:00-04:00 New_York
		ts,val
		2020-10-04T09:00:00-04:00 New_York,T
		2020-10-04T21:15:00-04:00 New_York,F
		2020-10-05T00:00:00-04:00 New_York,F
		`
	emptyRes string = `ver:"3.0"
		empty
		`
)
