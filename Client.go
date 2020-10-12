package haystack

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
)

// Client models a client connection to a server using the Haystack API.
type Client struct {
	httpClient  *http.Client
	uri         string
	username    string
	password    string
	authHeaders map[string]string
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

	return &Client{
		httpClient:  &http.Client{Timeout: timeout},
		uri:         uri,
		username:    username,
		password:    password,
		authHeaders: map[string]string{},
	}
}

// Open simply opens and authenticates the connection
func (client *Client) Open() error {
	helloAuth, helloErr := client.sendHello()
	if helloErr != nil {
		return helloErr
	}

	// TODO Expand support to multiple auth scheme returns
	// TODO Add support for non-scram auth
	if helloAuth.scheme != "scram" {
		return errors.New("Auth scheme not supported: " + helloAuth.scheme)
	}

	handshakeToken := helloAuth.get("handshakeToken")
	hashName := helloAuth.get("hash")

	scramErr := client.openScram(handshakeToken, hashName)
	return scramErr
}

// sendHello sends a hello message and returns the value of the WWW-Authenticate header
func (client *Client) sendHello() (authMsg, error) {
	req, _ := http.NewRequest("get", client.uri+"about", nil)
	reqAuth := authMsg{
		scheme: "hello",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(client.username)),
		},
	}
	client.prepare(req)
	req.Header.Add("Authorization", reqAuth.toString())

	// TODO delete me
	// fmt.Println(req.URL.String())
	// for key, vals := range req.Header {
	// 	fmt.Print(key)
	// 	fmt.Print(": ")
	// 	for _, val := range vals {
	// 		fmt.Print(val)
	// 		fmt.Print(", ")
	// 	}
	// 	fmt.Println(" ")
	// }

	resp, respErr := client.httpClient.Do(req)
	if respErr != nil {
		return authMsg{}, respErr
	}
	if resp.StatusCode != 401 {
		return authMsg{}, NewHTTPError(resp.StatusCode, "Hello "+resp.Status)
	}
	resp.Body.Close()
	respAuthString := resp.Header.Get("WWW-Authenticate")
	if respAuthString == "" {
		return authMsg{}, NewAuthError("Missing required header: WWW-Authenticate")
	}
	return authMsgFromString(respAuthString), nil
}

// openScram opens a scram connection. 'hashName' only supports "SHA-256" and "SHA-512"
func (client *Client) openScram(handshakeToken string, hashName string) error {
	var hash func() hash.Hash
	if hashName == "SHA-256" {
		hash = sha256.New
	} else if hashName == "SHA-512" {
		hash = sha512.New
	} else { // Only support SHA-256 and SHA-512
		return NewAuthError("Auth hash not supported: " + hashName)
	}

	var in []byte
	var scram = NewScram(hash, client.username, client.password)
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("get", client.uri+"about", nil)
		reqAuth := authMsg{
			scheme: "scram",
			attrs: map[string]string{
				"handshakeToken": handshakeToken,
				"data":           encoding.EncodeToString(out),
			},
		}
		client.prepare(req)
		req.Header.Add("Authorization", reqAuth.toString())
		resp, _ := client.httpClient.Do(req)

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusOK { // We expect unauthorized until complete.
			return NewHTTPError(resp.StatusCode, resp.Status)
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
		return scram.Err()
	}

	finalAuth := authMsg{
		scheme: "bearer",
		attrs: map[string]string{ // Only keep the authToken
			"authToken": authToken,
		},
	}
	client.authHeaders["Authorization"] = finalAuth.toString()

	return nil
}

// About calls the 'about' op.
func (client *Client) About() (Dict, error) {
	result, err := client.Call("about", EmptyGrid())
	if err != nil {
		return Dict{}, err
	}
	return result.RowAt(0).ToDict(), nil
}

// Ops calls the 'ops' op.
func (client *Client) Ops() (Grid, error) {
	return client.Call("ops", EmptyGrid())
}

// Formats calls the 'formats' op.
func (client *Client) Formats() (Grid, error) {
	return client.Call("formats", EmptyGrid())
}

// Read calls the 'read' op with a filter and no result limit.
func (client *Client) Read(filter string) (Grid, error) {
	return client.ReadLimit(filter, 0)
}

// ReadLimit calls the 'read' op with a filter and a result limit.
func (client *Client) ReadLimit(filter string, limit int) (Grid, error) {
	var limitVal Val
	if limit <= 0 {
		limitVal = NewNull()
	} else {
		limitVal = NewNumber(float64(limit), "")
	}
	var gb GridBuilder
	gb.AddColNoMeta("filter")
	gb.AddColNoMeta("limit")
	gb.AddRow([]Val{
		NewStr(filter),
		limitVal,
	})
	return client.Call("read", gb.ToGrid())
}

// ReadByIds calls the 'read' op with the input ids.
func (client *Client) ReadByIds(ids []Ref) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("id")
	for _, id := range ids {
		gb.AddRow([]Val{
			id,
		})
	}
	return client.Call("read", gb.ToGrid())
}

// Nav calls the 'nav' op to navigate a project for learning and discovery
func (client *Client) Nav(navId Val) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("navId")
	gb.AddRow([]Val{
		navId,
	})
	return client.Call("nav", gb.ToGrid())
}

// PointWriteArray calls the 'pointWrite' op to query the point write priority array for the input id.
func (client *Client) PointWriteArray(id Ref) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("id")
	gb.AddRow([]Val{
		id,
	})
	return client.Call("pointWrite", gb.ToGrid())
}

// PointWrite calls the 'pointWrite' op to write the val to the given point.
func (client *Client) PointWrite(id Ref, level int, val Val, who string, duration Number) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("id")
	gb.AddColNoMeta("level")
	gb.AddColNoMeta("val")
	gb.AddColNoMeta("who")
	gb.AddColNoMeta("duration")
	gb.AddRow([]Val{
		id,
		NewNumber(float64(level), ""),
		val,
		NewStr(who),
		duration,
	})
	return client.Call("pointWrite", gb.ToGrid())
}

// HisReadAbsDate calls the 'hisRead' op with an input absolute Date range.
func (client *Client) HisReadAbsDate(id Ref, from Date, to Date) (Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

// HisReadAbsDateTime calls the 'hisRead' op with an input absolute DateTime range.
func (client *Client) HisReadAbsDateTime(id Ref, from DateTime, to DateTime) (Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

// HisRead calls the 'hisRead' op with the given range string. See Haystack API docs for accepted rangeString values.
func (client *Client) HisRead(id Ref, rangeString string) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("id")
	gb.AddColNoMeta("range")
	gb.AddRow([]Val{
		id,
		NewStr(rangeString),
	})
	return client.Call("hisRead", gb.ToGrid())
}

// HisWrite calls the 'hisWrite' op with the given id and Dicts of history items. Only the "ts" and "val" fields from
// the history items are included.
func (client *Client) HisWrite(id Ref, hisItems []Dict) (Grid, error) {
	var gb GridBuilder
	gb.AddMetaVal("id", id)
	gb.AddColNoMeta("ts")
	gb.AddColNoMeta("val")
	gb.AddRowDicts(hisItems)
	return client.Call("hisWrite", gb.ToGrid())
}

// InvokeAction calls the 'invokeAction' op with the given id, action name, and arguments.
func (client *Client) InvokeAction(id Ref, action string, args map[string]Val) (Grid, error) {
	var gb GridBuilder
	gb.AddMetaVal("id", id)
	gb.AddMetaVal("action", NewStr(action))

	rowVals := []Val{}
	for name, val := range args {
		gb.AddColNoMeta(name)
		rowVals = append(rowVals, val)
	}
	gb.AddRow(rowVals)
	return client.Call("invokeAction", gb.ToGrid())
}

// Eval calls the 'eval' op to evaluate a vendor specific expression.
func (client *Client) Eval(expr string) (Grid, error) {
	var gb GridBuilder
	gb.AddColNoMeta("expr")
	gb.AddRow([]Val{NewStr(expr)})
	return client.Call("eval", gb.ToGrid())
}

// Call executes the given operation. The request grid is posted to the client URI and the response is parsed as a grid.
func (client *Client) Call(op string, reqGrid Grid) (Grid, error) {
	req := reqGrid.ToZinc()
	resp, err := client.postString(op, req)
	if err != nil {
		return EmptyGrid(), err
	}

	var reader ZincReader
	reader.InitString(resp)
	val := reader.ReadVal()
	switch val := val.(type) {
	case Grid:
		if val.Meta().Get("err") != NewNull() {
			return EmptyGrid(), NewCallError(val)
		}
		return val, nil
	default:
		return EmptyGrid(), errors.New("Result was not a grid")
	}
}

func (client *Client) postString(op string, reqBody string) (string, error) {
	reqReader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("post", client.uri+op, reqReader)
	client.prepare(req)
	req.Header.Add("Connection", "Close")
	resp, respErr := client.httpClient.Do(req)
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

func (client *Client) prepare(req *http.Request) {
	for name, val := range client.authHeaders {
		req.Header.Add(name, val)
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "text/zinc; charset=utf-8")
	req.Header.Add("Accept", "text/zinc")
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
		builder.WriteRune('=')
		builder.WriteString(val)
	}
	return builder.String()
}

/////////////////////
// Errors
/////////////////////

// AuthError represents an error in the authorization
type AuthError struct {
	message string
}

// NewAuthError creates a new AuthError object.
func NewAuthError(message string) AuthError {
	return AuthError{message: message}
}

func (err AuthError) Error() string {
	return "Auth error: " + err.message
}

// CallError occurs when communication is successful, and a Grid is returned, but the grid has an err
// tag indicating a server side error.
type CallError struct {
	grid Grid
}

// NewCallError creates a new CallError object.
func NewCallError(grid Grid) CallError {
	return CallError{grid: grid}
}

func (err CallError) Error() string {
	dis := err.grid.Meta().Get("dis")
	switch val := dis.(type) {
	case Str:
		return "Call error: " + val.String()
	default:
		return "Call error: Server side error"
	}
}

// HTTPError occurs when communication is successful with a server, but we receive an HTTP error response.
type HTTPError struct {
	code    int
	message string
}

// NewHTTPError creates a new HTTPError object.
func NewHTTPError(code int, message string) HTTPError {
	return HTTPError{code: code, message: message}
}

func (err HTTPError) Error() string {
	return "HTTP error: " + err.message
}

// NetworkError occurs when there is a network I/O or connection problem with communication to the server.
type NetworkError struct {
	message string
}

// NewNetworkError creates a new NetworkError object.
func NewNetworkError(message string) NetworkError {
	return NetworkError{message: message}
}

func (err NetworkError) Error() string {
	return "Network error: " + err.message
}
