package client

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"hash"
	"net/http"
	"strings"

	"github.com/NeedleInAJayStack/haystack/auth"
)

// scramAuthenticator

type scramAuthenticator struct {
	clientHTTP clientHTTP
	uri        string
	username   string
	password   string
	initialMsg authMsg
}

func (authenticator scramAuthenticator) authorizationHeader() (string, error) {
	hashName := authenticator.initialMsg.get("hash")
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
	var scram = auth.NewScram(hash, authenticator.username, authenticator.password)
	handshakeToken := authenticator.initialMsg.get("handshakeToken")
	var authToken string
	for !scram.Step(in) {
		out := scram.Out()

		req, _ := http.NewRequest("GET", authenticator.uri, nil)
		reqAuth := authMsg{
			scheme: "scram",
			attrs: map[string]string{
				"handshakeToken": handshakeToken,
				"data":           encoding.EncodeToString(out),
			},
		}
		setStandardHeaders(req, reqAuth.toString())
		resp, _ := authenticator.clientHTTP.do(req)
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

// plaintextAuthenticator

type plaintextAuthenticator struct {
	clientHTTP clientHTTP
	uri        string
	username   string
	password   string
}

func (authenticator plaintextAuthenticator) authorizationHeader() (string, error) {
	reqAuth := authMsg{
		scheme: "plaintext",
		attrs: map[string]string{
			"username": encoding.EncodeToString([]byte(authenticator.username)),
			"password": encoding.EncodeToString([]byte(authenticator.password)),
		},
	}
	req, _ := http.NewRequest("GET", authenticator.uri, nil)
	setStandardHeaders(req, reqAuth.toString())
	resp, _ := authenticator.clientHTTP.do(req)
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

// basicAuthenticator

type basicAuthenticator struct {
	clientHTTP clientHTTP
	uri        string
	username   string
	password   string
}

func (authenticator basicAuthenticator) authorizationHeader() (string, error) {
	authValue := base64.StdEncoding.EncodeToString([]byte(authenticator.username + ":" + authenticator.password))
	basicAuth := "Basic " + authValue

	// Test the basic auth to ensure that it works
	req, _ := http.NewRequest("GET", authenticator.uri, nil)
	setStandardHeaders(req, basicAuth)
	resp, _ := authenticator.clientHTTP.do(req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", NewAuthError("Basic auth failed with status: " + resp.Status)
	}

	return basicAuth, nil
}
