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

var encoding = base64.RawURLEncoding
var userAgent = "go"

// Open simply opens and authenticates the connection
func (client *Client) Open() error {
	helloAuth, helloErr := client.sendHello()
	if helloErr != nil {
		return helloErr
	}

	// TODO Expand support to multiple auth scheme returns
	scheme, helloAttrs := parseAuth(helloAuth)
	scheme = strings.ToLower(scheme)
	// TODO Add support for non-scram auth
	if scheme != "scram" {
		return errors.New("Auth scheme not supported: " + scheme)
	}

	handshakeToken := helloAttrs["handshakeToken"]
	hashName := helloAttrs["hash"]

	scramErr := client.openScram(handshakeToken, hashName)
	return scramErr
}

// sendHello sends a hello message and returns the value of the WWW-Authenticate header
func (client *Client) sendHello() (string, error) {
	// Send Hello message
	req, _ := http.NewRequest("get", client.uri+"about", nil)
	scheme := "hello"
	attrs := map[string]string{
		"username": encoding.EncodeToString([]byte(client.username)),
	}
	reqAuth := buildAuth(scheme, attrs)
	req.Header.Add("Authorization", reqAuth)
	resp, respErr := client.httpClient.Do(req)
	if respErr != nil {
		return "", respErr
	}
	if resp.StatusCode != 401 {
		return "", errors.New(resp.Status)
	}
	resp.Body.Close()
	auth := resp.Header.Get("WWW-Authenticate")
	if auth == "" {
		return "", errors.New("Missing required header: WWW-Authenticate")
	}
	return auth, nil
}

// openScram opens a scram connection. 'hashName' only supports "SHA-256" and "SHA-512"
func (client *Client) openScram(handshakeToken string, hashName string) error {
	var hash func() hash.Hash
	if hashName == "SHA-256" {
		hash = sha256.New
	} else if hashName == "SHA-512" {
		hash = sha512.New
	} else { // Only support SHA-256 and SHA-512
		return errors.New("Auth hash not supported: " + hashName)
	}

	var in []byte
	var scram = NewScram(hash, client.username, client.password)
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("get", client.uri+"about", nil)
		reqAttrs := map[string]string{
			"handshakeToken": handshakeToken,
			"data":           encoding.EncodeToString(out),
		}
		reqAuth := buildAuth("scram", reqAttrs)

		req.Header.Add("Authorization", reqAuth)
		resp, _ := client.httpClient.Do(req)

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusOK { // We expect unauthorized until complete.
			return errors.New(resp.Status)
		}
		respAuth := resp.Header.Get("WWW-Authenticate")
		if respAuth == "" { // This header switches to Authentication-Info on success
			respAuth = resp.Header.Get("Authentication-Info")
		}
		_, respAttrs := parseAuth(respAuth)

		handshakeToken = respAttrs["handshakeToken"] // it grows over time
		dataEnc := respAttrs["data"]
		authToken = respAttrs["authToken"] // This will only be set on the last message
		data, _ := encoding.DecodeString(dataEnc)

		in = data
	}
	if scram.Err() != nil {
		return scram.Err()
	}

	authAttrs := map[string]string{ // Only keep the authToken
		"authToken": authToken,
	}
	client.authHeaders["Authorization"] = buildAuth("bearer", authAttrs)

	return nil
}

// Parses an authentication message, and returns the scheme, and the attributes
// Assumes auth messages follow the form: "[scheme] <name1>=<val1>, <name2>=<val2>, ..."
func parseAuth(str string) (string, map[string]string) {
	attrs := make(map[string]string)
	attributeStrs := strings.Split(str, ",")
	scheme := ""

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

	return scheme, attrs
}

// Aggregates an authentication message and returns the result
// Resulting auth messages follow the form: "[scheme] <name1>=<val1>, <name2>=<val2>, ..."
func buildAuth(scheme string, attrs map[string]string) string {
	builder := new(strings.Builder)
	if scheme != "" {
		builder.WriteString(strings.ToUpper(scheme))
		builder.WriteRune(' ')
	}
	firstVal := true
	for name, val := range attrs {
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

func (client *Client) Call(op string, reqGrid haystack.Grid) (haystack.Grid, error) {
	req := reqGrid.ToZinc()
	resp, err := client.postString(op, req)
	if err != nil {
		return haystack.Grid{}, err
	}

	var reader io.ZincReader
	reader.InitString(resp)
	val := reader.ReadVal()
	switch val := val.(type) {
	case haystack.Grid:
		return val, nil
	default:
		return haystack.Grid{}, errors.New("Result was not a grid")
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
		return "", errors.New("http response: " + resp.Status)
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
