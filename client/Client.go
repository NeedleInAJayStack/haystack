package client

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"hash"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/NeedleInAJayStack/haystack"
	"github.com/NeedleInAJayStack/haystack/auth"
	"github.com/NeedleInAJayStack/haystack/io"
)

// Client models a client connection to a server using the Haystack API.
type Client struct {
	clientHTTP clientHTTP
	uri        string
	auth       string
}

var encoding = base64.RawURLEncoding
var userAgent = "Go-haystack-client"

// NewClient creates a new Client object.
func NewClient(uri string, username string, password string) *Client {
	// check URI
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		panic("URI isn't http or https: " + uri)
	}
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	timeout, _ := time.ParseDuration("1m")
	client := http.Client{Timeout: timeout}

	return &Client{
		clientHTTP: &clientHTTPImpl{
			uri:        uri,
			username:   username,
			password:   password,
			httpClient: &client,
		},
		uri:  uri,
		auth: "",
	}
}

// Open simply opens and authenticates the connection
func (client *Client) Open() error {
	auth, err := client.clientHTTP.getAuthHeader()
	if err != nil {
		return err
	}
	client.auth = auth
	return nil
}

// About calls the 'about' op.
func (client *Client) About() (haystack.Dict, error) {
	result, err := client.Call("about", haystack.EmptyGrid())
	if err != nil {
		return haystack.Dict{}, err
	}
	return result.RowAt(0).ToDict(), nil
}

// Close closes and de-authenticates the client
func (client *Client) Close() error {
	_, err := client.Call("close", haystack.EmptyGrid())
	return err
}

// Defs calls the 'defs' op.
func (client *Client) Defs() (haystack.Grid, error) {
	return client.Call("defs", haystack.EmptyGrid())
}

// DefsWithFilter calls the 'defs' op with a filter grid.
func (client *Client) DefsWithFilter(filter string, limit int) (haystack.Grid, error) {
	return client.Call("defs", filterGrid(filter, limit))
}

// Libs calls the 'libs' op.
func (client *Client) Libs() (haystack.Grid, error) {
	return client.Call("libs", haystack.EmptyGrid())
}

// LibsWithFilter calls the 'libs' op with a filter grid.
func (client *Client) LibsWithFilter(filter string, limit int) (haystack.Grid, error) {
	return client.Call("libs", filterGrid(filter, limit))
}

// Ops calls the 'ops' op.
func (client *Client) Ops() (haystack.Grid, error) {
	return client.Call("ops", haystack.EmptyGrid())
}

// OpsWithFilter calls the 'ops' op with a filter grid.
func (client *Client) OpsWithFilter(filter string, limit int) (haystack.Grid, error) {
	return client.Call("ops", filterGrid(filter, limit))
}

// Filetypes calls the 'filetypes' op.
func (client *Client) Filetypes() (haystack.Grid, error) {
	return client.Call("filetypes", haystack.EmptyGrid())
}

// FiletypesWithFilter calls the 'filetypes' op with a filter grid.
func (client *Client) FiletypesWithFilter(filter string, limit int) (haystack.Grid, error) {
	return client.Call("filetypes", filterGrid(filter, limit))
}

// Read calls the 'read' op with a filter and no result limit.
func (client *Client) Read(filter string) (haystack.Grid, error) {
	return client.ReadLimit(filter, 0)
}

// ReadLimit calls the 'read' op with a filter and a result limit.
func (client *Client) ReadLimit(filter string, limit int) (haystack.Grid, error) {
	return client.Call("read", filterGrid(filter, limit))
}

// ReadByIds calls the 'read' op with the input ids.
func (client *Client) ReadByIds(ids []haystack.Ref) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("id")
	for _, id := range ids {
		gb.AddRow([]haystack.Val{id})
	}
	return client.Call("read", gb.ToGrid())
}

// Nav calls the 'nav' op to navigate a project for learning and discovery
func (client *Client) Nav(navId haystack.Val) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("navId")
	gb.AddRow([]haystack.Val{navId})
	return client.Call("nav", gb.ToGrid())
}

// WatchSubCreate calls the 'watchSub' op to create a new subscription. If `lease` is 0 or less, no lease is added
// to the subscription
func (client *Client) WatchSubCreate(
	watchDis string,
	lease haystack.Number,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	meta := map[string]haystack.Val{"watchDis": haystack.NewStr(watchDis)}
	if lease.Float() > 0 {
		meta["lease"] = lease
	}

	gb := haystack.NewGridBuilder()
	gb.AddMeta(meta)
	gb.AddColNoMeta("ids")
	for _, id := range ids {
		gb.AddRow([]haystack.Val{id})
	}
	return client.Call("watchSub", gb.ToGrid())
}

// WatchSubAdd calls the 'watchSub' op to add to an existing subscription. If `lease` is 0 or less, no lease is added
// to the subscription.
func (client *Client) WatchSubAdd(
	watchId string,
	lease haystack.Number,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	meta := map[string]haystack.Val{"watchId": haystack.NewStr(watchId)}
	if lease.Float() > 0 {
		meta["lease"] = lease
	}

	gb := haystack.NewGridBuilder()
	gb.AddMeta(meta)
	gb.AddColNoMeta("ids")
	for _, id := range ids {
		gb.AddRow([]haystack.Val{id})
	}
	return client.Call("watchSub", gb.ToGrid())
}

// WatchUnsub calls the 'watchUnsub' op to delete or remove entities from a existing subscription. If `lease` is 0
// or less, no lease is added to the subscription.
func (client *Client) WatchUnsub(
	watchId string,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	meta := map[string]haystack.Val{"watchId": haystack.NewStr(watchId)}
	if len(ids) <= 0 {
		meta["close"] = haystack.NewMarker()
	}

	gb := haystack.NewGridBuilder()
	gb.AddMeta(meta)
	gb.AddColNoMeta("ids")
	for _, id := range ids {
		gb.AddRow([]haystack.Val{id})
	}
	return client.Call("watchUnsub", gb.ToGrid())
}

// WatchPoll calls the 'watchPoll' op to poll values of a subscription.
func (client *Client) WatchPoll(
	watchId string,
	refresh bool,
) (haystack.Grid, error) {
	meta := map[string]haystack.Val{"watchId": haystack.NewStr(watchId)}
	if refresh {
		meta["refresh"] = haystack.NewMarker()
	}

	gb := haystack.NewGridBuilder()
	gb.AddMeta(meta)
	return client.Call("watchPoll", gb.ToGrid())
}

// PointWriteStatus calls the 'pointWrite' op to query the point write priority array status for the input id.
func (client *Client) PointWriteStatus(id haystack.Ref) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("id")
	gb.AddRow([]haystack.Val{id})
	return client.Call("pointWrite", gb.ToGrid())
}

// PointWrite calls the 'pointWrite' op to write the val to the given point.
func (client *Client) PointWrite(
	id haystack.Ref,
	level int,
	val haystack.Val,
	who string,
	duration haystack.Number,
) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("id")
	gb.AddColNoMeta("level")
	gb.AddColNoMeta("val")
	gb.AddColNoMeta("who")
	gb.AddColNoMeta("duration")
	gb.AddRow([]haystack.Val{
		id,
		haystack.NewNumber(float64(level), ""),
		val,
		haystack.NewStr(who),
		duration,
	})
	return client.Call("pointWrite", gb.ToGrid())
}

// HisReadAbsDate calls the 'hisRead' op with an input absolute Date range.
func (client *Client) HisReadAbsDate(id haystack.Ref, from haystack.Date, to haystack.Date) (haystack.Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

// HisReadAbsDateTime calls the 'hisRead' op with an input absolute DateTime range.
func (client *Client) HisReadAbsDateTime(id haystack.Ref, from haystack.DateTime, to haystack.DateTime) (haystack.Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

// HisRead calls the 'hisRead' op with the given range string. See Haystack API docs for accepted rangeString values.
func (client *Client) HisRead(id haystack.Ref, rangeString string) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("id")
	gb.AddColNoMeta("range")
	gb.AddRow([]haystack.Val{
		id,
		haystack.NewStr(rangeString),
	})
	return client.Call("hisRead", gb.ToGrid())
}

// HisWrite calls the 'hisWrite' op with the given id and Dicts of history items. Only the "ts" and "val" fields from
// the history items are included.
func (client *Client) HisWrite(id haystack.Ref, hisItems []haystack.Dict) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddMetaVal("id", id)
	gb.AddColNoMeta("ts")
	gb.AddColNoMeta("val")
	gb.AddRowDicts(hisItems)
	return client.Call("hisWrite", gb.ToGrid())
}

// InvokeAction calls the 'invokeAction' op with the given id, action name, and arguments.
func (client *Client) InvokeAction(id haystack.Ref, action string, args map[string]haystack.Val) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddMetaVal("id", id)
	gb.AddMetaVal("action", haystack.NewStr(action))

	rowVals := []haystack.Val{}
	for name, val := range args {
		gb.AddColNoMeta(name)
		rowVals = append(rowVals, val)
	}
	gb.AddRow(rowVals)
	return client.Call("invokeAction", gb.ToGrid())
}

// Eval calls the 'eval' op to evaluate a vendor specific expression.
func (client *Client) Eval(expr string) (haystack.Grid, error) {
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("expr")
	gb.AddRow([]haystack.Val{haystack.NewStr(expr)})
	return client.Call("eval", gb.ToGrid())
}

// Call executes the given operation. The request grid is posted to the client URI and the response is parsed as a grid.
func (client *Client) Call(op string, reqGrid haystack.Grid) (haystack.Grid, error) {
	req := reqGrid.ToZinc()

	resp, err := client.clientHTTP.postString(client.uri, client.auth, op, req)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	var reader io.ZincReader
	reader.InitString(resp)
	val, err := reader.ReadVal()
	if err != nil {
		return haystack.EmptyGrid(), err
	}
	switch val := val.(type) {
	case haystack.Grid:
		if val.Meta().Get("err") != haystack.NewNull() {
			return haystack.EmptyGrid(), NewCallError(val)
		}
		return val, nil
	default:
		return haystack.EmptyGrid(), errors.New("result was not a grid")
	}
}

// filterGrid creates a Grid consisting of a `filter` Str and `limit` Number columns.
// If a value of 0 or less is passed to limit, no limit is applied.
func filterGrid(filter string, limit int) haystack.Grid {
	var limitVal haystack.Val
	if limit <= 0 {
		limitVal = haystack.NewNull()
	} else {
		limitVal = haystack.NewNumber(float64(limit), "")
	}
	gb := haystack.NewGridBuilder()
	gb.AddColNoMeta("filter")
	gb.AddColNoMeta("limit")
	gb.AddRow([]haystack.Val{
		haystack.NewStr(filter),
		limitVal,
	})
	return gb.ToGrid()
}

// authMsg models a message in the Haystack authorization header format.
// They follow the form: "[scheme] <name1>=<val1>, <name2>=<val2>, ..."
type authMsg struct {
	scheme string
	attrs  map[string]string
}

func authMsgFromString(str string) authMsg {
	attrs := make(map[string]string)
	attributeStrs := strings.Split(str, ",")
	scheme := ""

	// The first one MAY include the scheme but not necessarily. Handle both situations
	firstAttr := attributeStrs[0]
	if strings.Contains(attributeStrs[0], " ") {
		schemeSplit := strings.Split(firstAttr, " ")
		scheme = strings.TrimSpace(schemeSplit[0])
		attributeStrs[0] = schemeSplit[1]
	}

	for _, attributeStr := range attributeStrs {
		attributeSplit := strings.Split(attributeStr, "=")
		name := strings.TrimSpace(attributeSplit[0])
		val := strings.TrimSpace(attributeSplit[1])
		attrs[name] = val
	}

	return authMsg{
		scheme: scheme,
		attrs:  attrs,
	}
}

func (authMsg *authMsg) get(attrName string) string {
	return authMsg.attrs[attrName]
}

func (authMsg *authMsg) toString() string {
	builder := new(strings.Builder)
	if authMsg.scheme != "" {
		builder.WriteString(strings.ToUpper(authMsg.scheme))
		builder.WriteRune(' ')
	}
	firstVal := true
	for name, val := range authMsg.attrs {
		if firstVal {
			firstVal = false
		} else {
			builder.WriteString(", ")
		}
		builder.WriteString(name)
		if val != "" {
			builder.WriteRune('=')
			builder.WriteString(val)
		}
	}
	return builder.String()
}

// clientHTTP is defined as an interface to allow dependency-injection testing
type clientHTTP interface {
	// getAuthHeader returns the `Authorization` header to use
	getAuthHeader() (string, error)
	// postString posts the given request body to the given URI and returns the response body
	postString(uri string, auth string, op string, reqBody string) (string, error)
}

// clientHTTPImpl is the default implementation of clientHTTP
type clientHTTPImpl struct {
	uri        string
	username   string
	password   string
	httpClient *http.Client
}

func (clientHTTP *clientHTTPImpl) getAuthHeader() (string, error) {
	req, _ := http.NewRequest("GET", clientHTTP.uri+"about", nil)
	reqAuth := authMsg{
		scheme: "hello",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(clientHTTP.username)),
		},
	}
	clientHTTP.setStandardHeaders(req, reqAuth.toString())

	resp, respErr := clientHTTP.httpClient.Do(req)
	if respErr != nil {
		return "", respErr
	}
	if resp.StatusCode != 401 {
		return "", NewHTTPError(resp.StatusCode, "`about` endpoint with HELLO scheme returned a non 401 status: "+resp.Status)
	}
	resp.Body.Close()
	respAuthString := resp.Header.Get("WWW-Authenticate")
	if respAuthString == "" {
		return "", NewAuthError("Missing required header: WWW-Authenticate")
	}
	helloAuth := authMsgFromString(respAuthString)

	var authToken string
	var authErr error
	switch strings.ToUpper(helloAuth.scheme) {
	case "SCRAM":
		authToken, authErr = clientHTTP.authTokenFromScram(helloAuth.get("handshakeToken"), helloAuth.get("hash"))
	case "PLAINTEXT":
		authToken, authErr = clientHTTP.authTokenFromPlaintext()
	default:
		return "", NewAuthError("Auth scheme not supported: " + helloAuth.scheme)
	}
	if authErr != nil {
		return "", authErr
	}

	finalAuth := authMsg{
		scheme: "bearer",
		attrs: map[string]string{
			"authToken": authToken,
		},
	}
	return finalAuth.toString(), nil
}

func (clientHTTP *clientHTTPImpl) postString(uri string, auth string, op string, reqBody string) (string, error) {
	reqReader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("POST", uri+op, reqReader)
	clientHTTP.setStandardHeaders(req, auth)
	req.Header.Add("Connection", "Close")
	resp, respErr := clientHTTP.httpClient.Do(req)
	if respErr != nil {
		return "", respErr
	}
	if resp.StatusCode != http.StatusOK {
		return "", NewHTTPError(resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return string(body), err
}

func (clientHTTP *clientHTTPImpl) authTokenFromScram(
	uri string,
	username string,
	password string,
	handshakeToken string,
	hashName string,
) (string, error) {
	var hash func() hash.Hash
	switch strings.ToUpper(hashName) {
	case "SHA-256":
		hash = sha256.New
	case "SHA-512":
		hash = sha512.New
	default: // Only support SHA-256 and SHA-512
		return "", NewAuthError("Auth hash not supported: " + hashName)
	}

	var in []byte
	var scram = auth.NewScram(hash, clientHTTP.username, clientHTTP.password)
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("GET", clientHTTP.uri+"about", nil)
		reqAuth := authMsg{
			scheme: "scram",
			attrs: map[string]string{
				"handshakeToken": handshakeToken,
				"data":           encoding.EncodeToString(out),
			},
		}
		clientHTTP.setStandardHeaders(req, reqAuth.toString())
		resp, _ := clientHTTP.httpClient.Do(req)

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusOK { // We expect unauthorized until complete.
			return "", NewHTTPError(resp.StatusCode, resp.Status)
		}
		respAuthString := resp.Header.Get("WWW-Authenticate")
		if respAuthString == "" { // This header switches to Authentication-Info on success
			respAuthString = resp.Header.Get("Authentication-Info")
		}
		respAuth := authMsgFromString(respAuthString)

		handshakeToken = respAuth.get("handshakeToken") // it grows over time
		dataEnc := respAuth.get("data")
		authToken = respAuth.get("authToken") // This will only be set on the last message
		data, _ := encoding.DecodeString(dataEnc)

		in = data
	}
	if scram.Err() != nil {
		return "", scram.Err()
	}
	return authToken, nil
}

func (clientHTTP *clientHTTPImpl) authTokenFromPlaintext() (string, error) {
	reqAuth := authMsg{
		scheme: "plaintext",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(clientHTTP.username)),
			"password": encoding.EncodeToString([]byte(clientHTTP.password)),
		},
	}
	req, _ := http.NewRequest("GET", clientHTTP.uri+"about", nil)
	clientHTTP.setStandardHeaders(req, reqAuth.toString())
	resp, _ := clientHTTP.httpClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return "", NewHTTPError(resp.StatusCode, resp.Status)
	}
	respAuthString := resp.Header.Get("Authentication-Info")
	respAuth := authMsgFromString(respAuthString)
	return respAuth.get("authToken"), nil
}

func (clientHTTP *clientHTTPImpl) setStandardHeaders(req *http.Request, auth string) {
	req.Header.Add("Authorization", auth)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "text/zinc; charset=utf-8")
	req.Header.Add("Accept", "text/zinc")
}

/////////////////////
// Errors
/////////////////////

// AuthError represents an error in the authorization
type AuthError struct {
	Message string
}

// NewAuthError creates a new AuthError object.
func NewAuthError(message string) AuthError {
	return AuthError{Message: message}
}

func (err AuthError) Error() string {
	return "Auth error: " + err.Message
}

// CallError occurs when communication is successful, and a Grid is returned, but the grid has an err
// tag indicating a server side error.
type CallError struct {
	Grid haystack.Grid
}

// NewCallError creates a new CallError object.
func NewCallError(grid haystack.Grid) CallError {
	return CallError{Grid: grid}
}

func (err CallError) Error() string {
	dis := err.Grid.Meta().Get("dis")
	switch val := dis.(type) {
	case haystack.Str:
		return "Call error: " + val.String()
	default:
		return "Call error: Server side error"
	}
}

// HTTPError occurs when communication is successful with a server, but we receive an HTTP error response.
type HTTPError struct {
	Code    int
	Message string
}

// NewHTTPError creates a new HTTPError object.
func NewHTTPError(code int, message string) HTTPError {
	return HTTPError{Code: code, Message: message}
}

func (err HTTPError) Error() string {
	return "HTTP error: " + err.Message
}

// NetworkError occurs when there is a network I/O or connection problem with communication to the server.
type NetworkError struct {
	Message string
}

// NewNetworkError creates a new NetworkError object.
func NewNetworkError(message string) NetworkError {
	return NetworkError{Message: message}
}

func (err NetworkError) Error() string {
	return "Network error: " + err.Message
}
