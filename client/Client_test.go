package client

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/NeedleInAJayStack/haystack"
	haystackIO "github.com/NeedleInAJayStack/haystack/io"
	"github.com/stretchr/testify/assert"
)

func TestClient_Open(t *testing.T) {
	client := testPostClient()
	openErr := client.Open()
	if openErr != nil {
		t.Error(openErr)
	}
}

func TestClient_About(t *testing.T) {
	client := testPostClient()
	actual, aboutErr := client.About()
	assert.Nil(t, aboutErr)
	// About returns a Dict so get the first value of expected manually
	expectedZinc := clientHTTPMock_about
	var reader haystackIO.ZincReader
	reader.InitString(expectedZinc)
	expectedGrid, err := reader.ReadVal()
	assert.Nil(t, err)
	expected := expectedGrid.(haystack.Grid).RowAt(0).ToDict()

	assert.Equal(t, actual, expected)
}

func TestClient_Close(t *testing.T) {
	client := testPostClient()
	err := client.Close()
	assert.Nil(t, err)
}

func TestClient_Filetypes(t *testing.T) {
	actual, err := testPostClient().Filetypes()
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_filetypes, t)

	get, err := testGetClient().Filetypes()
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_filetypes, t)
}

func TestClient_Ops(t *testing.T) {
	actual, err := testPostClient().Ops()
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_ops, t)

	get, err := testGetClient().Ops()
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_ops, t)
}

func TestClient_Read(t *testing.T) {
	actual, err := testPostClient().Read("site")
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_readSites, t)

	get, err := testGetClient().Read("site")
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_readSites, t)
}

func TestClient_ReadLimit(t *testing.T) {
	actual, err := testPostClient().ReadLimit("point", 1)
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_readPoint, t)

	get, err := testGetClient().ReadLimit("point", 1)
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_readPoint, t)
}

func TestClient_ReadByIds(t *testing.T) {
	points, pointsErr := testPostClient().ReadLimit("point", 1)
	assert.Nil(t, pointsErr)
	pointRef := points.RowAt(0).Get("id").(haystack.Ref)

	actual, err := testPostClient().ReadByIds([]haystack.Ref{pointRef})
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_readPoint, t)

	get, err := testGetClient().ReadByIds([]haystack.Ref{pointRef})
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_readPoint, t)
}

func TestClient_HisRead(t *testing.T) {
	points, pointsErr := testPostClient().ReadLimit("point", 1)
	assert.Nil(t, pointsErr)
	pointRef := points.RowAt(0).Get("id").(haystack.Ref)

	actual, err := testPostClient().HisRead(pointRef, "yesterday")
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_hisRead20210103, t)

	get, err := testGetClient().HisRead(pointRef, "yesterday")
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_hisRead20210103, t)
}

func TestClient_HisReadAbsDate(t *testing.T) {
	points, pointsErr := testPostClient().ReadLimit("point", 1)
	assert.Nil(t, pointsErr)
	pointRef := points.RowAt(0).Get("id").(haystack.Ref)

	fromDate := haystack.NewDate(2020, 10, 4)
	toDate := haystack.NewDate(2020, 10, 5)

	actual, err := testPostClient().HisReadAbsDate(pointRef, fromDate, toDate)
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_hisRead20201004to6, t)

	get, err := testGetClient().HisReadAbsDate(pointRef, fromDate, toDate)
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_hisRead20201004to6, t)
}

func TestClient_HisReadAbsDateTime(t *testing.T) {
	points, pointsErr := testPostClient().ReadLimit("point", 1)
	assert.Nil(t, pointsErr)
	pointRef := points.RowAt(0).Get("id").(haystack.Ref)

	fromTs, _ := haystack.NewDateTimeFromString("2020-10-04T00:00:00-07:00 Los_Angeles")
	toTs, _ := haystack.NewDateTimeFromString("2020-10-05T00:00:00-07:00 Los_Angeles")

	actual, err := testPostClient().HisReadAbsDateTime(pointRef, fromTs, toTs)
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_hisReadDateTimes, t)

	get, err := testGetClient().HisReadAbsDateTime(pointRef, fromTs, toTs)
	assert.Nil(t, err)
	testClient_ValZinc(get, clientHTTPMock_hisReadDateTimes, t)
}

func TestClient_WatchSubCreate(t *testing.T) {
	actual, err := testPostClient().WatchSubCreate(
		"abc",
		haystack.NewNumber(1, "min"),
		[]haystack.Ref{haystack.NewRef("abc-123", "")},
	)
	assert.Nil(t, err)
	testClient_ValZinc(actual, emptyRes, t)
}

func TestClient_WatchSubAdd(t *testing.T) {
	actual, err := testPostClient().WatchSubAdd(
		"abc",
		haystack.NewNumber(1, "min"),
		[]haystack.Ref{haystack.NewRef("abc-123", "")},
	)
	assert.Nil(t, err)
	testClient_ValZinc(actual, emptyRes, t)
}

func TestClient_WatchUnsub(t *testing.T) {
	actual, err := testPostClient().WatchUnsub("abc", []haystack.Ref{haystack.NewRef("abc-123", "")})
	assert.Nil(t, err)
	testClient_ValZinc(actual, emptyRes, t)
}

func TestClient_Eval(t *testing.T) {
	actual, err := testPostClient().Eval("read(point)")
	assert.Nil(t, err)
	testClient_ValZinc(actual, clientHTTPMock_readPoint, t)
}

func testClient_ValZinc(actual haystack.Val, expectedZinc string, t *testing.T) {
	var reader haystackIO.ZincReader
	reader.InitString(expectedZinc)
	expected, err := reader.ReadVal()
	assert.Nil(t, err)
	assert.Equal(t, actual, expected, "\nACTUAL:\n"+actual.ToZinc()+"\n\nEXPECT:\n"+expected.ToZinc())
}

func testPostClient() *Client {
	client := NewClient(
		"http://localhost:8080/api/demo/",
		"test",
		"test",
		NoAuthenticator{},
	)
	client.clientHTTP = &clientHTTPMock{}
	return client
}

func testGetClient() *Client {
	return &Client{
		clientHTTP: &clientHTTPMock{},
		method:     Get,
		uri:        "http://localhost:8080/api/demo/",
		username:   "test",
		password:   "test",
	}
}

// clientHTTPMock allows us to remove the HTTP dependency within tests
type clientHTTPMock struct{}

func (clientHTTPMock *clientHTTPMock) do(req *http.Request) (*http.Response, error) {
	response := http.Response{
		Header: make(http.Header),
		Body:   http.NoBody,
	}
	response.StatusCode = 500
	var responseBody string
	var err error
	switch req.Method {
	case "GET":
		urlSlice := strings.Split(req.URL.Path, "/")
		op := urlSlice[len(urlSlice)-1]
		params := map[string]string{}
		for name, values := range req.URL.Query() {
			params[name] = values[0]
		}
		responseBody, err = clientHTTPMock.getResponse(op, params)
	case "POST":
		urlSlice := strings.Split(req.URL.Path, "/")
		op := urlSlice[len(urlSlice)-1]
		reqBody, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			return &response, readErr
		}
		responseBody, err = clientHTTPMock.postResponse(op, string(reqBody))
	}
	if err != nil {
		return &response, err
	}
	response.StatusCode = 200
	response.Body = io.NopCloser(strings.NewReader(responseBody))
	return &response, nil
}

func (clientHTTPMock *clientHTTPMock) postResponse(op string, reqBody string) (string, error) {
	// These are taken from a SkySpark 3.0.26 demo project on 2021-01-03
	switch op {
	case "about":
		// Can't use string literal because of Uri backticks
		return clientHTTPMock_about, nil
	case "close":
		return emptyRes, nil
	case "filetypes":
		return clientHTTPMock_filetypes, nil
	case "ops":
		return clientHTTPMock_ops, nil
	case "read":
		if reqBody == "ver:\"3.0\"\nfilter, limit\n\"site\", N" { // readAll sites
			return clientHTTPMock_readSites, nil
		} else if reqBody == "ver:\"3.0\"\nfilter, limit\n\"point\", 1" { // readLimit point
			return clientHTTPMock_readPoint, nil
		} else if reqBody == "ver:\"3.0\"\nid\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" { // readById
			return clientHTTPMock_readPoint, nil
		}
		return emptyRes, errors.New("'read' argument not supported by mock class")
	case "hisRead":
		if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"yesterday\"" { // hisRead relative
			return clientHTTPMock_hisRead20210103, nil
		} else if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"2020-10-04,2020-10-05\"" { // hisRead absolute dates
			return clientHTTPMock_hisRead20201004to6, nil
		} else if reqBody == "ver:\"3.0\"\nid, range\n@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\", \"2020-10-04T00:00:00-07:00 Los_Angeles,2020-10-05T00:00:00-07:00 Los_Angeles\"" { // hisRead absolute datetimes
			return clientHTTPMock_hisReadDateTimes, nil
		}
		return emptyRes, errors.New("'hisRead' argument not supported by mock class")
	case "watchSub":
		return emptyRes, nil
	case "watchUnsub":
		return emptyRes, nil
	case "eval":
		if reqBody == "ver:\"3.0\"\nexpr\n\"read(point)\"" { // eval read(point)
			return clientHTTPMock_readPoint, nil
		}
		return emptyRes, errors.New("'eval' argument not supported by mock class")
	default:
		return emptyRes, errors.New("haystack op not supported by mock class: " + op)
	}
}

func (clientHTTPMock *clientHTTPMock) getResponse(op string, params map[string]string) (string, error) {
	// These are taken from a SkySpark 3.0.26 demo project on 2021-01-03
	switch op {
	case "hello":
		return emptyRes, nil
	case "about":
		return clientHTTPMock_about, nil
	case "close":
		return emptyRes, nil
	case "filetypes":
		return clientHTTPMock_filetypes, nil
	case "ops":
		return clientHTTPMock_ops, nil
	case "read":
		if params["filter"] == "\"site\"" && params["limit"] == "N" { // readAll sites
			return clientHTTPMock_readSites, nil
		} else if params["filter"] == "\"point\"" && params["limit"] == "1" { // readLimit point
			return clientHTTPMock_readPoint, nil
		} else if params["id"] == "@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" { // readById
			return clientHTTPMock_readPoint, nil
		}
		return emptyRes, errors.New("'read' argument not supported by mock class")
	case "hisRead":
		if params["id"] == "@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" && params["range"] == "\"yesterday\"" { // hisRead relative
			return clientHTTPMock_hisRead20210103, nil
		} else if params["id"] == "@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" && params["range"] == "\"2020-10-04,2020-10-05\"" { // hisRead absolute dates
			return clientHTTPMock_hisRead20201004to6, nil
		} else if params["id"] == "@p:demo:r:2725da26-1dda68ee \"Gaithersburg RTU-1 Fan\"" && params["range"] == "\"2020-10-04T00:00:00-07:00 Los_Angeles,2020-10-05T00:00:00-07:00 Los_Angeles\"" { // hisRead absolute datetimes
			return clientHTTPMock_hisReadDateTimes, nil
		}
		return emptyRes, errors.New("'hisRead' argument not supported by mock class")
	case "watchSub":
		return emptyRes, nil
	case "watchUnsub":
		return emptyRes, nil
	case "eval":
		return emptyRes, nil
	default:
		return emptyRes, errors.New("haystack op not supported by mock class: " + op)
	}
}

const (
	clientHTTPMock_about string = "ver:\"3.0\"\n" + // Can't use string literal because of Uri backticks
		"haystackVersion,projName,serverName,serverBootTime,serverTime,productName,productUri,productVersion,moduleName,moduleVersion,tz,whoami,hostDis,hostModel,hostId\n" +
		"\"3.0\",\"demo\",\"JaysDesktop\",2021-01-03T00:21:01.588-07:00 Denver,2021-01-03T00:21:43.799-07:00 Denver,\"SkySpark\",`http://skyfoundry.com/skyspark`,\"3.0.26\",\"skyarcd\",\"3.0.26\",\"Denver\",\"test\",\"Linux amd64 5.4.0-58-generic\",\"Linux amd64 5.4.0-58-generic\",NA\n"
	clientHTTPMock_filetypes string = `ver:"3.0"
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
		def, doc, is, lib, noSideEffects, typeName
		^op:invokeAction, "Invoke a user action on a target entity.\nSee docHaystack::Ops#invokeAction chapter.", [^op], ^lib:phIoT, N, N
		^op:hisWrite, "Write historized time series data from a his-point.\nSee docHaystack::Ops#hisWrite chapter.", [^op], ^lib:phIoT, N, "hx::HxHisWriteOp"
		^op:watchUnsub, "Unsubscribe to entity data.\nSee docHaystack::Ops#watchUnsub chapter.", [^op], ^lib:ph, N, "hx::HxWatchUnsubOp"
		^op:read, "Query the a set of entity records by id or by filter.\nSee docHaystack::Ops#read chapter.", [^op], ^lib:ph, M, "hx::HxReadOp"
		^op:ops, "Query the op defs in the current namespace.\nSee docHaystack::Ops#ops chapter.", [^op], ^lib:ph, M, "hx::HxOpsOp"
		^op:close, "Close the current session and cancel the auth bearer token.\nSee docHaystack::Ops#close chapter.", [^op], ^lib:ph, N, "hx::HxCloseOp"
		^op:defs, "Query the definitions in the current namespace.\nSee docHaystack::Ops#defs chapter.", [^op], ^lib:ph, M, "hx::HxDefsOp"
		^op:hisRead, "Read historized time series data from a his-point.\nSee docHaystack::Ops#hisRead chapter.", [^op], ^lib:phIoT, M, "hx::HxHisReadOp"
		^op:watchSub, "Subscribe to entity data.\nSee docHaystack::Ops#watchSub chapter.", [^op], ^lib:ph, N, "hx::HxWatchSubOp"
		^op:pointWrite, "Read or command a writable-point.\nSee docHaystack::Ops#pointWrite chapter.", [^op], ^lib:phIoT, N, "hx::HxPointWriteOp"
		^op:eval, "Evaluate an Axon expression", [^op], ^lib:hx, N, "hx::HxEvalOp"
		^op:libs, "Query the lib defs in the current namespace.\nSee docHaystack::Ops#libs chapter.", [^op], ^lib:ph, M, "hx::HxLibsOp"
		^op:watchPoll, "Poll a watch subscription.\nSee docHaystack::Ops#watchPoll chapter.", [^op], ^lib:ph, N, "hx::HxWatchPollOp"
		^op:filetypes, "Query the filetype defs in the current namespace.\nSee docHaystack::Ops#filetypes chapter.", [^op], ^lib:ph, M, "hx::HxFiletypesOp"
		^op:commit, "Commit one or more diffs to the Folio database", [^op], ^lib:hx, N, "hx::HxCommitOp"
		^op:about, "Query basic information about the server.\nSee docHaystack::Ops#about chapter.", [^op], ^lib:ph, M, "hx::HxAboutOp"
		^op:nav, "Query the navigation tree for discovery.\nSee docHaystack::Ops#nav chapter.", [^op], ^lib:ph, M, "hx::HxNavOp"
		`
	clientHTTPMock_readSites string = `ver:"3.0"
		id,area,dis,geoAddr,geoCity,geoCoord,geoCountry,geoPostalCode,geoState,geoStreet,hq,metro,occupiedEnd,occupiedStart,primaryFunction,regionRef,site,store,storeNum,tz,weatherStationRef,yearBuilt,mod
		@p:demo:r:2725da26-ac563571 "Headquarters",140797ft²,"Headquarters","600 W Main St, Richmond, VA","Richmond",C(37.545826,-77.449188),"US","23220","VA","600 W Main St",M,"Richmond",18:00:00,08:00:00,"Office",@p:demo:r:2725da26-f3e488bc "Richmond",M,,,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1999,2020-10-23T18:15:02.701Z
		@p:demo:r:2725da26-3ca6125c "Gaithersburg",8013ft²,"Gaithersburg","18212 Montgomery Village Ave, Gaithersburg, MD","Gaithersburg",C(39.154824,-77.209002),"US","20879","MD","18212 Montgomery Village Ave",,"Washington DC",21:00:00,09:00:00,"Retail Store",@p:demo:r:2725da26-e77a16f1 "Washington DC",M,M,4,"New_York",@p:demo:r:2725da26-9bb170b8 "Washington, DC",2001,2020-10-23T18:15:02.797Z
		@p:demo:r:2725da26-d280b1b5 "Short Pump",17122ft²,"Short Pump","11282 W Broad St, Richmond, VA","Glen Allen",C(37.650338,-77.606105),"US","23060","VA","11282 W Broad St",,"Richmond",21:00:00,10:00:00,"Retail Store",@p:demo:r:2725da26-f3e488bc "Richmond",M,M,3,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1999,2020-10-23T18:15:02.763Z
		@p:demo:r:2725da26-505b4ae8 "Carytown",3149ft²,"Carytown","3504 W Cary St, Richmond, VA","Richmond",C(37.555385,-77.486903),"US","23221","VA","3504 W Cary St",,"Richmond",20:00:00,10:00:00,"Retail Store",@p:demo:r:2725da26-f3e488bc "Richmond",M,M,1,"New_York",@p:demo:r:2725da26-9fd27896 "Richmond, VA",1996,2020-10-23T18:15:02.742Z
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
