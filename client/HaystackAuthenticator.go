package client

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"net/http"
	"strings"

	"github.com/NeedleInAJayStack/haystack/auth"
)

// HaystackAuthenticator is an Authenticator that uses the standard Haystack authentication schemes
// See https://project-haystack.org/doc/Auth for more detail
type HaystackAuthenticator struct{}

func (authenticator HaystackAuthenticator) Authenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
) (map[string]string, error) {
	aboutUri := uri + "about"
	req, _ := http.NewRequest("GET", aboutUri, nil)
	reqAuth := authMsg{
		scheme: "hello",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(username)),
		},
	}
	setHeaders(req, map[string]string{"Authorization": reqAuth.toString()})

	resp, respErr := client.do(req)
	if respErr != nil {
		return map[string]string{}, respErr
	}
	// If we get 200, authentication is not required
	if resp.StatusCode == 200 {
		return map[string]string{}, nil
	}
	respWwwAuthenticate := resp.Header.Get("WWW-Authenticate")
	resp.Body.Close()
	if resp.StatusCode != 401 {
		return map[string]string{}, NewHTTPError(resp.StatusCode, "`about` endpoint with HELLO scheme returned a non 401 status: "+resp.Status)
	}
	if respWwwAuthenticate == "" {
		return map[string]string{}, NewAuthError("No WWW-Authenticate header in response")
	}

	helloAuth := authMsgFromString(respWwwAuthenticate)

	var authHeader string
	var authErr error
	switch strings.ToUpper(helloAuth.scheme) {
	case "SCRAM":
		authHeader, authErr = authenticator.scramAuthenticate(
			aboutUri,
			username,
			password,
			client,
			helloAuth,
		)
	case "PLAINTEXT":
		authHeader, authErr = authenticator.plaintextAuthenticate(
			aboutUri,
			username,
			password,
			client,
		)
	default:
		return map[string]string{}, NewAuthError("Auth scheme not supported: " + helloAuth.scheme)
	}
	if authErr != nil {
		return map[string]string{}, authErr
	}
	return map[string]string{"Authorization": authHeader}, nil
}

func (authenticator HaystackAuthenticator) scramAuthenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
	initialMsg authMsg,
) (string, error) {
	hashName := initialMsg.get("hash")
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
	handshakeToken := initialMsg.get("handshakeToken")
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("GET", uri, nil)
		reqAuth := authMsg{
			scheme: "scram",
			attrs: map[string]string{
				"handshakeToken": handshakeToken,
				"data":           encoding.EncodeToString(out),
			},
		}
		setHeaders(req, map[string]string{"Authorization": reqAuth.toString()})
		resp, _ := client.do(req)
		defer resp.Body.Close()

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

	finalAuth := authMsg{
		scheme: "bearer",
		attrs: map[string]string{
			"authToken": authToken,
		},
	}
	return finalAuth.toString(), nil
}

func (authenticator HaystackAuthenticator) plaintextAuthenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
) (string, error) {
	reqAuth := authMsg{
		scheme: "plaintext",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(username)),
			"password": encoding.EncodeToString([]byte(password)),
		},
	}
	req, _ := http.NewRequest("GET", uri, nil)
	setHeaders(req, map[string]string{"Authorization": reqAuth.toString()})
	resp, _ := client.do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewHTTPError(resp.StatusCode, resp.Status)
	}
	respAuthString := resp.Header.Get("Authentication-Info")
	respAuth := authMsgFromString(respAuthString)
	authToken := respAuth.get("authToken")

	finalAuth := authMsg{
		scheme: "bearer",
		attrs: map[string]string{
			"authToken": authToken,
		},
	}
	return finalAuth.toString(), nil
}
