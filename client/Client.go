package client

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/NeedleInAJayStack/haystack"
	"github.com/NeedleInAJayStack/haystack/io"
)

// Client models a client connection to a server using the Haystack API.
type Client struct {
	clientHTTP clientHTTP
	uri        string
	username   string
	password   string
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
	jar, _ := cookiejar.New(nil) // Required to persist client cookies

	return &Client{
		clientHTTP: &clientHTTPImpl{
			&http.Client{
				Timeout: timeout,
				Jar:     jar,
			},
		},
		uri:      uri,
		username: username,
		password: password,
		auth:     "",
	}
}

// Open simply opens and authenticates the connection
func (client *Client) Open() error {
	auth, err := client.getAuthHeader()
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

	resp, err := client.postString(client.uri, client.auth, op, req)
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

// getAuthHeader returns the `Authorization` header to use
func (client *Client) getAuthHeader() (string, error) {
	req, _ := http.NewRequest("GET", client.authUri(), nil)
	reqAuth := authMsg{
		scheme: "hello",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(client.username)),
		},
	}
	setStandardHeaders(req, reqAuth.toString())

	resp, respErr := client.clientHTTP.do(req)
	if respErr != nil {
		return "", respErr
	}
	// If we get 200, authentication is not required
	if resp.StatusCode == 200 {
		return "", nil
	}
	respWwwAuthenticate := resp.Header.Get("WWW-Authenticate")
	respServer := resp.Header.Get("Server")
	respSetCookie := resp.Header.Get("Set-Cookie")
	resp.Body.Close()
	if resp.StatusCode != 401 {
		return "", NewHTTPError(resp.StatusCode, "`about` endpoint with HELLO scheme returned a non 401 status: "+resp.Status)
	}

	var authErr error

	// First try Haystack standard authentication scheme
	if respWwwAuthenticate != "" {
		haystackAuthHeader, haystackErr := client.haystackAuth(respWwwAuthenticate)
		if haystackErr == nil {
			return haystackAuthHeader, nil
		} else {
			authErr = haystackErr
		}
	}

	// If we can't authenticate with Haystack, try basic auth
	isBasicAuth := strings.Contains(strings.ToLower(respWwwAuthenticate), "basic")
	isNiagara := strings.Contains(strings.ToLower(respServer), "niagara") || strings.Contains(strings.ToLower(respSetCookie), "niagara")
	if isBasicAuth || isNiagara {
		basicAuthHeader, basicErr := client.basicAuthenticator().authorizationHeader()
		if basicErr == nil {
			return basicAuthHeader, nil
		} else {
			authErr = basicErr
		}
	}

	if authErr == nil {
		authErr = NewAuthError("No suitable auth scheme found")
	}
	return "", authErr
}

func (client *Client) haystackAuth(wwwAuthenticate string) (string, error) {
	helloAuth := authMsgFromString(wwwAuthenticate)

	var authHeader string
	var authErr error
	switch strings.ToUpper(helloAuth.scheme) {
	case "SCRAM":
		authHeader, authErr = client.scramAuthenticator(helloAuth).authorizationHeader()
	case "PLAINTEXT":
		authHeader, authErr = client.plaintextAuthenticator().authorizationHeader()
	default:
		return "", NewAuthError("Auth scheme not supported: " + helloAuth.scheme)
	}
	if authErr != nil {
		return "", authErr
	}
	return authHeader, nil
}

func (client *Client) scramAuthenticator(initialMsg authMsg) scramAuthenticator {
	return scramAuthenticator{
		clientHTTP: client.clientHTTP,
		uri:        client.authUri(),
		username:   client.username,
		password:   client.password,
		initialMsg: initialMsg,
	}
}

func (client *Client) plaintextAuthenticator() plaintextAuthenticator {
	return plaintextAuthenticator{
		clientHTTP: client.clientHTTP,
		uri:        client.authUri(),
		username:   client.username,
		password:   client.password,
	}
}

func (client *Client) basicAuthenticator() basicAuthenticator {
	return basicAuthenticator{
		clientHTTP: client.clientHTTP,
		uri:        client.authUri(),
		username:   client.username,
		password:   client.password,
	}
}

func (client *Client) postString(uri string, auth string, op string, reqBody string) (string, error) {
	reqReader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("POST", uri+op, reqReader)
	setStandardHeaders(req, auth)
	req.Header.Add("Connection", "Close")
	resp, respErr := client.clientHTTP.do(req)
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

// Returns the URL used to authenticate the client
func (client *Client) authUri() string {
	return client.uri + "about"
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

func setStandardHeaders(req *http.Request, auth string) {
	req.Header.Add("Authorization", auth)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "text/zinc; charset=utf-8")
	req.Header.Add("Accept", "text/zinc")
}
