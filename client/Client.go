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

	"gitlab.com/NeedleInAJayStack/haystack"
	"gitlab.com/NeedleInAJayStack/haystack/io"
)

type Client struct {
	httpClient  *http.Client
	uri         string
	username    string
	password    string
	authHeaders map[string]string
}

var encoding = base64.RawURLEncoding
var userAgent = "go"

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

func (client *Client) IsOpen() bool {
	return client.authHeaders["Authorization"] == ""
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
	req.Header.Add("Authorization", reqAuth.toString())
	resp, respErr := client.httpClient.Do(req)
	if respErr != nil {
		return authMsg{}, respErr
	}
	if resp.StatusCode != 401 {
		return authMsg{}, NewHttpError(resp.StatusCode, resp.Status)
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

		req.Header.Add("Authorization", reqAuth.toString())
		resp, _ := client.httpClient.Do(req)

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusOK { // We expect unauthorized until complete.
			return NewHttpError(resp.StatusCode, resp.Status)
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

func (client *Client) About() (haystack.Dict, error) {
	result, err := client.Call("about", haystack.EmptyGrid())
	if err != nil {
		return haystack.Dict{}, err
	}
	return result.RowAt(0).ToDict(), nil
}

func (client *Client) Ops() (haystack.Grid, error) {
	return client.Call("ops", haystack.EmptyGrid())
}

func (client *Client) Formats() (haystack.Grid, error) {
	return client.Call("formats", haystack.EmptyGrid())
}

func (client *Client) Read(filter string) (haystack.Grid, error) {
	return client.ReadLimit(filter, 0)
}

func (client *Client) ReadLimit(filter string, limit int) (haystack.Grid, error) {
	var limitVal haystack.Val
	if limit <= 0 {
		limitVal = haystack.NewNull()
	} else {
		limitVal = haystack.NewNumber(float64(limit), "")
	}
	var gb haystack.GridBuilder
	gb.AddColNoMeta("filter")
	gb.AddColNoMeta("limit")
	gb.AddRow([]haystack.Val{
		haystack.NewStr(filter),
		limitVal,
	})
	return client.Call("read", gb.ToGrid())
}

func (client *Client) HisReadAbsDate(id haystack.Ref, from haystack.Date, to haystack.Date) (haystack.Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

func (client *Client) HisReadAbsDateTime(id haystack.Ref, from haystack.DateTime, to haystack.DateTime) (haystack.Grid, error) {
	rangeString := from.ToZinc() + "," + to.ToZinc()
	return client.HisRead(id, rangeString)
}

func (client *Client) HisRead(id haystack.Ref, rangeString string) (haystack.Grid, error) {
	var gb haystack.GridBuilder
	gb.AddColNoMeta("id")
	gb.AddColNoMeta("range")
	gb.AddRow([]haystack.Val{
		id,
		haystack.NewStr(rangeString),
	})
	return client.Call("hisRead", gb.ToGrid())
}

func (client *Client) Eval(expr string) (haystack.Grid, error) {
	var gb haystack.GridBuilder
	gb.AddColNoMeta("expr")
	gb.AddRow([]haystack.Val{haystack.NewStr(expr)})
	return client.Call("eval", gb.ToGrid())
}

func (client *Client) Call(op string, reqGrid haystack.Grid) (haystack.Grid, error) {
	req := reqGrid.ToZinc()
	resp, err := client.postString(op, req)
	if err != nil {
		return haystack.EmptyGrid(), err
	}

	var reader io.ZincReader
	reader.InitString(resp)
	val := reader.ReadVal()
	switch val := val.(type) {
	case haystack.Grid:
		if val.Meta().Get("err") != haystack.NewNull() {
			return haystack.EmptyGrid(), NewCallError(val)
		} else {
			return val, nil
		}
	default:
		return haystack.EmptyGrid(), errors.New("Result was not a grid")
	}
}

func (client *Client) postString(op string, reqBody string) (string, error) {
	reqReader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("post", client.uri+op, reqReader)
	client.prepare(req)
	req.Header.Add("Connection", "Close")
	req.Header.Add("Content-Type", "text/zinc; charset=utf-8") // TODO support more mimeTypes beyond UTF-8 zinc
	resp, respErr := client.httpClient.Do(req)
	if respErr != nil {
		return "", respErr
	}
	if resp.StatusCode != http.StatusOK {
		return "", NewHttpError(resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func (client *Client) prepare(req *http.Request) {
	for name, val := range client.authHeaders {
		req.Header.Add(name, val)
	}
}
