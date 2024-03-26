package client

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/NeedleInAJayStack/haystack/auth"
)

// clientHTTP is defined as an interface to allow dependency-injection testing
type clientHTTP interface {
	// getAuthHeader returns the `Authorization` header to use
	getAuthHeader(uri string, username string, password string) (string, error)
	// postString posts the given request body to the given URI and returns the response body
	postString(uri string, auth string, op string, reqBody string) (string, error)
}

// clientHTTPImpl is the default implementation of clientHTTP
type clientHTTPImpl struct {
	httpClient *http.Client
}

func (clientHTTP *clientHTTPImpl) getAuthHeader(uri string, username string, password string) (string, error) {
	req, _ := http.NewRequest("GET", uri+"about", nil)
	reqAuth := authMsg{
		scheme: "hello",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(username)),
		},
	}
	clientHTTP.setStandardHeaders(req, reqAuth.toString())

	resp, respErr := clientHTTP.httpClient.Do(req)
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
	if respWwwAuthenticate == "" {
		return "", NewAuthError("Missing required header: WWW-Authenticate")
	}
	if resp.StatusCode != 401 {
		return "", NewHTTPError(resp.StatusCode, "`about` endpoint with HELLO scheme returned a non 401 status: "+resp.Status)
	}

	// First try Haystack standard authentication scheme
	haystackAuthHeader, haystackErr := clientHTTP.haystackAuth(uri, username, password, respWwwAuthenticate)
	if haystackErr == nil {
		return haystackAuthHeader, nil
	}

	// If we can't authenticate with Haystack, try basic auth
	isBasicAuth := strings.Contains(strings.ToLower(respWwwAuthenticate), "basic")
	isNiagara := strings.Contains(strings.ToLower(respServer), "niagara") || strings.Contains(strings.ToLower(respSetCookie), "niagara")
	if isBasicAuth || isNiagara {
		return clientHTTP.basicAuth(uri, username, password)
	}

	return haystackAuthHeader, haystackErr
}

func (clientHTTP *clientHTTPImpl) haystackAuth(uri string, username string, password string, wwwAuthenticate string) (string, error) {
	helloAuth := authMsgFromString(wwwAuthenticate)

	var authToken string
	var authErr error
	switch strings.ToUpper(helloAuth.scheme) {
	case "SCRAM":
		authToken, authErr = clientHTTP.authTokenFromScram(uri, username, password, helloAuth.get("handshakeToken"), helloAuth.get("hash"))
	case "PLAINTEXT":
		authToken, authErr = clientHTTP.authTokenFromPlaintext(uri, username, password)
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

func (clientHTTP *clientHTTPImpl) basicAuth(uri string, username string, password string) (string, error) {
	authValue := username + ":" + password
	basicAuth := authMsg{
		scheme: "basic",
		attrs: map[string]string{
			authValue: "",
		},
	}

	// Test the basic auth to ensure that it works
	req, _ := http.NewRequest("GET", uri+"about", nil)
	clientHTTP.setStandardHeaders(req, basicAuth.toString())
	resp, _ := clientHTTP.httpClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return "", NewAuthError("Basic auth failed with status: " + resp.Status)
	}

	return basicAuth.toString(), nil
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
	var scram = auth.NewScram(hash, username, password)
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("GET", uri+"about", nil)
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

func (clientHTTP *clientHTTPImpl) authTokenFromPlaintext(
	uri string,
	username string,
	password string,
) (string, error) {
	reqAuth := authMsg{
		scheme: "plaintext",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(username)),
			"password": encoding.EncodeToString([]byte(password)),
		},
	}
	req, _ := http.NewRequest("GET", uri+"about", nil)
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
