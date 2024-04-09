package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/NeedleInAJayStack/haystack"
	haystackIO "github.com/NeedleInAJayStack/haystack/io"
)

// Client models a client connection to a server using the Haystack API.
type Client struct {
	method     ClientMethod
	uri        string
	username   string
	password   string
	clientHTTP ClientHTTP
	// Authenticator to use for authentication
	authenticator Authenticator
	// Headers managed by the authenticator
	authHeaders map[string]string
	// Headers that should be used in all requests. These may override the authentication headers.
	Headers map[string]string
}

var encoding = base64.RawURLEncoding
var userAgent = "Go-haystack-client"

// NewClient creates a new Client object.
func NewClient(uri string, username string, password string, authenticator Authenticator) *Client {
	// check URI
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		panic("URI isn't http or https: " + uri)
	}
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	timeout, _ := time.ParseDuration("1m")

	clientHTTP := &clientHTTPImpl{
		&http.Client{
			Timeout: timeout,
		},
	}

	return &Client{
		clientHTTP:    clientHTTP,
		uri:           uri,
		username:      username,
		password:      password,
		authenticator: authenticator,
		authHeaders:   map[string]string{},
		Headers:       map[string]string{},
	}
}

// Open simply opens and authenticates the connection
func (client *Client) Open() error {
	authHeaders, err := client.authenticator.Authenticate(client.uri, client.username, client.password, client.clientHTTP)
	if err != nil {
		return err
	}
	client.authHeaders = authHeaders
	return nil
}

// About calls the 'about' op.
func (client *Client) About() (haystack.Dict, error) {
	var result haystack.Grid
	var err error
	switch client.method {
	case Get:
		result, err = client.get("about", map[string]haystack.Val{})
	default:
		result, err = client.post("about", haystack.EmptyGrid())
	}
	if err != nil {
		return haystack.Dict{}, err
	}
	return result.RowAt(0).ToDict(), nil
}

// Close closes and de-authenticates the client
func (client *Client) Close() error {
	var err error
	switch client.method {
	case Get:
		return errors.New("'close' op does not support GET method")
	default:
		_, err = client.post("close", haystack.EmptyGrid())
	}
	client.authHeaders = map[string]string{}
	return err
}

// Defs calls the 'defs' op.
func (client *Client) Defs() (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("about", map[string]haystack.Val{})
	default:
		return client.post("about", haystack.EmptyGrid())
	}
}

// DefsWithFilter calls the 'defs' op with a filter grid.
func (client *Client) DefsWithFilter(filter string, limit int) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("defs", filterParams(filter, limit))
	default:
		return client.post("defs", filterGrid(filter, limit))
	}
}

// Libs calls the 'libs' op.
func (client *Client) Libs() (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("libs", map[string]haystack.Val{})
	default:
		return client.post("libs", haystack.EmptyGrid())
	}
}

// LibsWithFilter calls the 'libs' op with a filter grid.
func (client *Client) LibsWithFilter(filter string, limit int) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("libs", filterParams(filter, limit))
	default:
		return client.post("libs", filterGrid(filter, limit))
	}
}

// Ops calls the 'ops' op.
func (client *Client) Ops() (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("ops", map[string]haystack.Val{})
	default:
		return client.post("ops", haystack.EmptyGrid())
	}
}

// OpsWithFilter calls the 'ops' op with a filter grid.
func (client *Client) OpsWithFilter(filter string, limit int) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("ops", filterParams(filter, limit))
	default:
		return client.post("ops", filterGrid(filter, limit))
	}
}

// Filetypes calls the 'filetypes' op.
func (client *Client) Filetypes() (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("filetypes", map[string]haystack.Val{})
	default:
		return client.post("filetypes", haystack.EmptyGrid())
	}
}

// FiletypesWithFilter calls the 'filetypes' op with a filter grid.
func (client *Client) FiletypesWithFilter(filter string, limit int) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("filetypes", filterParams(filter, limit))
	default:
		return client.post("filetypes", filterGrid(filter, limit))
	}
}

// Read calls the 'read' op with a filter and no result limit.
func (client *Client) Read(filter string) (haystack.Grid, error) {
	return client.ReadLimit(filter, 0)
}

// ReadLimit calls the 'read' op with a filter and a result limit.
func (client *Client) ReadLimit(filter string, limit int) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("read", filterParams(filter, limit))
	default:
		return client.post("read", filterGrid(filter, limit))
	}
}

// ReadByIds calls the 'read' op with the input ids.
func (client *Client) ReadByIds(ids []haystack.Ref) (haystack.Grid, error) {
	switch client.method {
	case Get:
		if len(ids) > 1 {
			return haystack.EmptyGrid(), errors.New("'read' op only supports single-ref requests on GET method")
		}
		if len(ids) == 0 {
			return haystack.EmptyGrid(), nil
		}
		return client.get("read", map[string]haystack.Val{"id": ids[0]})
	default:
		gb := haystack.NewGridBuilder()
		gb.AddColNoMeta("id")
		for _, id := range ids {
			gb.AddRow([]haystack.Val{id})
		}
		return client.post("read", gb.ToGrid())
	}
}

// Nav calls the 'nav' op to navigate a project for learning and discovery
func (client *Client) Nav(navId haystack.Val) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return client.get("nav", map[string]haystack.Val{"navId": navId})
	default:
		gb := haystack.NewGridBuilder()
		gb.AddColNoMeta("navId")
		gb.AddRow([]haystack.Val{navId})
		return client.post("nav", gb.ToGrid())
	}
}

// WatchSubCreate calls the 'watchSub' op to create a new subscription. If `lease` is 0 or less, no lease is added
// to the subscription
func (client *Client) WatchSubCreate(
	watchDis string,
	lease haystack.Number,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'watchSub' op does not support GET method")
	default:
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
		return client.post("watchSub", gb.ToGrid())
	}
}

// WatchSubAdd calls the 'watchSub' op to add to an existing subscription. If `lease` is 0 or less, no lease is added
// to the subscription.
func (client *Client) WatchSubAdd(
	watchId string,
	lease haystack.Number,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'watchSub' op does not support GET method")
	default:
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
		return client.post("watchSub", gb.ToGrid())
	}
}

// WatchUnsub calls the 'watchUnsub' op to delete or remove entities from a existing subscription. If `lease` is 0
// or less, no lease is added to the subscription.
func (client *Client) WatchUnsub(
	watchId string,
	ids []haystack.Ref,
) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'watchUnsub' op does not support GET method")
	default:
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
		return client.post("watchUnsub", gb.ToGrid())
	}
}

// WatchPoll calls the 'watchPoll' op to poll values of a subscription.
func (client *Client) WatchPoll(
	watchId string,
	refresh bool,
) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'watchPoll' op does not support GET method")
	default:
		meta := map[string]haystack.Val{"watchId": haystack.NewStr(watchId)}
		if refresh {
			meta["refresh"] = haystack.NewMarker()
		}

		gb := haystack.NewGridBuilder()
		gb.AddMeta(meta)
		return client.post("watchPoll", gb.ToGrid())
	}
}

// PointWriteStatus calls the 'pointWrite' op to query the point write priority array status for the input id.
func (client *Client) PointWriteStatus(id haystack.Ref) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'pointWrite' op does not support GET method")
	default:
		gb := haystack.NewGridBuilder()
		gb.AddColNoMeta("id")
		gb.AddRow([]haystack.Val{id})
		return client.post("pointWrite", gb.ToGrid())
	}
}

// PointWrite calls the 'pointWrite' op to write the val to the given point.
func (client *Client) PointWrite(
	id haystack.Ref,
	level int,
	val haystack.Val,
	who string,
	duration haystack.Number,
) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'pointWrite' op does not support GET method")
	default:
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
		return client.post("pointWrite", gb.ToGrid())
	}
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
	switch client.method {
	case Get:
		return client.get("hisRead", map[string]haystack.Val{"id": id, "range": haystack.NewStr(rangeString)})
	default:
		gb := haystack.NewGridBuilder()
		gb.AddColNoMeta("id")
		gb.AddColNoMeta("range")
		gb.AddRow([]haystack.Val{
			id,
			haystack.NewStr(rangeString),
		})
		return client.post("hisRead", gb.ToGrid())
	}
}

// HisWrite calls the 'hisWrite' op with the given id and Dicts of history items. Only the "ts" and "val" fields from
// the history items are included.
func (client *Client) HisWrite(id haystack.Ref, hisItems []haystack.Dict) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'hisWrite' op does not support GET method")
	default:
		gb := haystack.NewGridBuilder()
		gb.AddMetaVal("id", id)
		gb.AddColNoMeta("ts")
		gb.AddColNoMeta("val")
		gb.AddRowDicts(hisItems)
		return client.post("hisWrite", gb.ToGrid())
	}
}

// InvokeAction calls the 'invokeAction' op with the given id, action name, and arguments.
func (client *Client) InvokeAction(id haystack.Ref, action string, args map[string]haystack.Val) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'invokeAction' op does not support GET method")
	default:
		gb := haystack.NewGridBuilder()
		gb.AddMetaVal("id", id)
		gb.AddMetaVal("action", haystack.NewStr(action))

		rowVals := []haystack.Val{}
		for name, val := range args {
			gb.AddColNoMeta(name)
			rowVals = append(rowVals, val)
		}
		gb.AddRow(rowVals)
		return client.post("invokeAction", gb.ToGrid())
	}
}

// Eval calls the 'eval' op to evaluate a vendor specific expression.
func (client *Client) Eval(expr string) (haystack.Grid, error) {
	switch client.method {
	case Get:
		return haystack.EmptyGrid(), errors.New("'eval' op does not support GET method")
	default:
		gb := haystack.NewGridBuilder()
		gb.AddColNoMeta("expr")
		gb.AddRow([]haystack.Val{haystack.NewStr(expr)})
		return client.post("eval", gb.ToGrid())
	}
}

// post executes the given operation. The request grid is posted to the client URI and the response is parsed as a grid.
func (client *Client) post(op string, reqGrid haystack.Grid) (haystack.Grid, error) {
	reqBody := reqGrid.ToZinc()

	reqReader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("POST", client.uri+op, reqReader)
	reqHeaders := client.authHeaders
	for key, value := range client.Headers {
		reqHeaders[key] = value
	}
	setHeaders(req, reqHeaders)
	req.Header.Add("Connection", "Close")
	resp, err := client.clientHTTP.do(req)
	if err != nil {
		return haystack.EmptyGrid(), err
	}
	if resp.StatusCode != http.StatusOK {
		return haystack.EmptyGrid(), NewHTTPError(resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	var reader haystackIO.ZincReader
	reader.InitString(string(respBody))
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

// post executes the given operation. The request grid is posted to the client URI and the response is parsed as a grid.
func (client *Client) get(op string, params map[string]haystack.Val) (haystack.Grid, error) {
	url := client.uri + op
	paramList := []string{}
	for name, val := range params {
		paramList = append(paramList, fmt.Sprintf("%v=%v", name, val.ToZinc()))
	}
	paramString := strings.Join(paramList, "&")
	if len(paramList) > 0 {
		url = url + "?" + paramString
	}

	req, _ := http.NewRequest("GET", url, strings.NewReader(""))
	reqHeaders := client.authHeaders
	for key, value := range client.Headers {
		reqHeaders[key] = value
	}
	setHeaders(req, reqHeaders)
	req.Header.Add("Connection", "Close")
	resp, err := client.clientHTTP.do(req)
	if err != nil {
		return haystack.EmptyGrid(), err
	}
	if resp.StatusCode != http.StatusOK {
		return haystack.EmptyGrid(), NewHTTPError(resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	var reader haystackIO.ZincReader
	reader.InitString(string(respBody))
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

// filterGrid creates a Grid consisting of a `filter` Str and `limit` Number columns.
// If a value of 0 or less is passed to limit, no limit is applied.
func filterParams(filter string, limit int) map[string]haystack.Val {
	var limitVal haystack.Val
	if limit <= 0 {
		limitVal = haystack.NewNull()
	} else {
		limitVal = haystack.NewNumber(float64(limit), "")
	}
	return map[string]haystack.Val{
		"filter": haystack.NewStr(filter),
		"limit":  limitVal,
	}
}

func setHeaders(req *http.Request, headers map[string]string) {
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "text/zinc; charset=utf-8")
	req.Header.Add("Accept", "text/zinc")

	for key, value := range headers {
		req.Header.Add(key, value)
	}
}
