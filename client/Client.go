package client

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

type Client struct {
	httpClient     *http.Client
	uri            string
	username       string
	password       string
	connectTimeout int
	readTimeout    int
}

func NewClient(uri string, username string, password string) *Client {
	// check URI
	if !strings.HasPrefix(uri, "http://") || !strings.HasPrefix(uri, "https://") {
		panic("URI isn't http or https: " + uri)
	}
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}

	return &Client{
		httpClient:     &http.Client{},
		uri:            uri,
		username:       username,
		password:       password,
		connectTimeout: 60 * 1000, // in milliseconds
		readTimeout:    60 * 1000, // in milliseconds
	}

}

// Open simply opens and authenticates the connection
func (client *Client) Open(uri string, username string, password string) error {
	helloRes, helloErr := client.sendHello()
	if helloErr != nil {
		return helloErr
	}

	// Attempt standard authentication via Haystack/RFC 7235
	openErr := client.openStd(helloRes)
	if openErr != nil {
		return openErr
	}

	// TODO PLACEHOLDER
	return nil
}

func (client *Client) sendHello() (*http.Response, error) {
	// Send Hello message
	username := []byte(client.username)
	req, _ := http.NewRequest("get", client.uri+"about", nil)
	scheme := "hello"
	attrs := map[string]string{
		"username": base64.RawURLEncoding.EncodeToString(username),
	}
	req.Header.Add("Authorization", buildAuth(scheme, attrs))
	return client.httpClient.Do(req)
}

func (client *Client) openStd(res *http.Response) error {
	if res.StatusCode != 401 {
		return errors.New(res.Status)
	}
	helloAuth := res.Header.Get("WWW-Authenticate")
	if helloAuth == "" {
		return errors.New("Missing required header: WWW-Authenticate")
	}

	// TODO Expand support to multiple auth scheme returns
	scheme, helloAttrs := parseAuth(helloAuth)
	scheme = strings.ToLower(scheme)
	// TODO Add support for non-scram auth
	if scheme != "scram" {
		return errors.New("Auth scheme not supported: " + scheme)
	}

	// Do SCRAM auth
	hash := helloAttrs["hash"]
	handshakeToken := helloAttrs["handshakeToken"]

	nonce := genNonce()

	req, _ := http.NewRequest("get", client.uri+"about", nil)
	initAttrs := map[string]string{
		"handshakeToken": handshakeToken,
	}
	req.Header.Add("Authorization", buildAuth(scheme, initAttrs))
	initReq, _ := http.NewRequest("get", client.uri+"about", nil)

	// TODO PLACEHOLDER
	return nil
}

// Parses an authentication message, and returns the scheme, and the attributes
// Assumes auth messages follow the form: "<scheme> <name1>=<val1>, <name2>=<val2>, ..."
func parseAuth(str string) (string, map[string]string) {
	firstSpaceIndex := strings.Index(str, " ")
	scheme := str[0:firstSpaceIndex]
	scheme = strings.ToLower(scheme)

	allAttrs := str[firstSpaceIndex+1 : len(str)]
	attributeStrs := strings.Split(allAttrs, ",")

	attrs := make(map[string]string)
	for _, attributeStr := range attributeStrs {
		attributeSplit := strings.Split(attributeStr, "=")
		name := attributeSplit[0]
		val := attributeSplit[1]
		attrs[name] = val
	}

	return scheme, attrs
}

func buildAuth(scheme string, attrs map[string]string) string {
	builder := new(strings.Builder)
	builder.WriteString(strings.ToUpper(scheme))
	builder.WriteRune(' ')
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

func genNonce() string {
	nonceSize := 16
	nonce := make([]byte, nonceSize)
	_, err := rand.Read(nonce) // This replaces the array values with random bytes
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(nonce)
}
