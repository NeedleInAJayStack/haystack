package client

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
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
func (client *Client) Open() error {
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
	req, _ := http.NewRequest("get", client.uri+"about", nil)
	scheme := "hello"
	attrs := map[string]string{
		"username": base64.RawURLEncoding.EncodeToString([]byte(client.username)),
	}
	req.Header.Add("Authorization", buildAuth(scheme, attrs))
	return client.httpClient.Do(req)
}

func (client *Client) openStd(helloRes *http.Response) error {
	if helloRes.StatusCode != 401 {
		return errors.New(helloRes.Status)
	}
	helloAuth := helloRes.Header.Get("WWW-Authenticate")
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

	handshakeToken := helloAttrs["handshakeToken"]
	hashFunc := helloAttrs["hash"]

	if hashFunc != "SHA-256" { // Only support SHA-256
		return errors.New("Auth hash not supported: " + hashFunc)
	}

	// Do SCRAM auth
	var in []byte
	var scram = NewScram(sha256.New, client.username, client.password)
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("get", client.uri+"about", nil)
		reqAttrs := map[string]string{
			"handshakeToken": handshakeToken,
			"data":           base64.RawURLEncoding.EncodeToString(out),
		}
		reqAuth := buildAuth(scheme, reqAttrs)

		// TODO DELETE ME. Debugging...
		fmt.Println("C: " + reqAuth)
		fmt.Println("    " + string(out))

		req.Header.Add("Authorization", reqAuth)
		res, _ := client.httpClient.Do(req)
		if res.StatusCode != 401 && res.StatusCode != 200 { // 401 is expected auth challenge, 200 is success
			return errors.New(res.Status)
		}
		resAuth := res.Header.Get("WWW-Authenticate")

		_, resAttrs := parseAuth(resAuth)

		handshakeToken = resAttrs["handshakeToken"] // it grows over time
		dataEnc := resAttrs["data"]
		data, _ := base64.RawURLEncoding.DecodeString(dataEnc)

		// TODO DELETE ME. Debugging...
		fmt.Println("S: " + resAuth)
		fmt.Println("    " + string(data))
		
		in = data
	}
	if scram.Err() != nil {
		return scram.Err()
	}

	// This was the work done before I found an existing implementation

	// nonce := "r=" + genNonce()
	// username := "n=" + client.username
	// clientFirstMessageBare := client.username + "," + nonce
	// clientFirstMessage := gs2Header() + clientFirstMessageBare

	// initReq, _ := http.NewRequest("get", client.uri+"about", nil)
	// initAttrs := map[string]string{
	// 	"handshakeToken": helloHandshakeToken,
	// 	"data":           base64.RawURLEncoding.EncodeToString([]byte(initMsg)),
	// }
	// initReq.Header.Add("Authorization", buildAuth(scheme, initAttrs))
	// initRes, _ := client.httpClient.Do(initReq)

	// initAuth := initRes.Header.Get("WWW-Authenticate")
	// scheme, initAttrs := parseAuth(initAuth)

	// initHandshakeToken := initAttrs["handshakeToken"]
	// initHashFunc := initAttrs["hash"]
	// initDataEnc := initAttrs["data"]
	// initDataStr, _ := string(base64.RawURLEncoding.DecodeString(initDataEnc))

	// initData := extractScramData(initDataStr)

	// finalNonceVal := initData["r"]
	// salt := initData["s"]
	// // iterationCount := initData["i"]

	// nonce = "r=" + finalNonceVal
	// cbindInput := gs2Header()
	// channelBinding := "c=" + base64.RawURLEncoding.EncodeToString(cbindInput)
	// clientFinalMessageWithoutProof := channelBinding + "," + nonce

	// // TODO Add support for other hash functions
	// if initHashFunc=="SHA-256"{
	// 	hash := sha256.New()
	// 	hash.Write([]byte(clientFinalMessageWithoutProof))

	// }
	// proof :=

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
		name := strings.TrimSpace(attributeSplit[0])
		val := strings.TrimSpace(attributeSplit[1])
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

func gs2Header() string {
	return "n,,"
}

// Extracts data from an X=ABCD,Y=1234 format to a map
func extractScramData(data string) map[string]string {
	dataParts := strings.Split(data, ",")
	result := make(map[string]string)
	for _, part := range dataParts {
		name := string(part[0])
		val := string(part[2:len(part)])
		result[name] = val
	}
	return result
}
