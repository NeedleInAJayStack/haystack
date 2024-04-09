package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHaystackAuthenticator_Plaintext(t *testing.T) {
	headers, err := HaystackAuthenticator{}.Authenticate(
		"http://localhost:8080/api/demo/",
		"test",
		"test",
		&clientHTTPPlaintextAuth{},
	)

	assert.Equal(t, headers, map[string]string{"Authorization": "BEARER authToken=pretend-this-is-a-token"})
	assert.Nil(t, err)
}

// clientHTTPPlaintextAuth validates the plaintext haystack authentication
// https://project-haystack.org/doc/docHaystack/Auth#plaintext
type clientHTTPPlaintextAuth struct{}

func (clientHTTPPlaintextAuth *clientHTTPPlaintextAuth) do(req *http.Request) (*http.Response, error) {
	response := http.Response{
		Header: make(http.Header),
		Body:   http.NoBody,
	}
	switch req.Method {
	case "GET":
		authMsg := authMsgFromString(req.Header.Get("Authorization"))
		if authMsg.scheme == "Bearer" {
			response.StatusCode = 200
		} else if authMsg.scheme == "PLAINTEXT" && authMsg.attrs["username"] == "dGVzdA" && authMsg.attrs["password"] == "dGVzdA" {
			response.StatusCode = 200
			response.Header.Set("Authentication-Info", "authToken=pretend-this-is-a-token")
		} else {
			response.StatusCode = 401
			response.Header.Set("WWW-Authenticate", "PLAINTEXT realm=\"Haystack\"")
		}
		return &response, nil
	}
	return &response, nil
}

// TODO: SCRAM
