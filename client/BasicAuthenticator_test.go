package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthenticator(t *testing.T) {
	headers, err := BasicAuthenticator{}.Authenticate(
		"http://localhost:8080/api/demo/",
		"test",
		"test",
		&clientHTTPBasicAuth{},
	)

	assert.Equal(t, headers, map[string]string{"Authorization": "Basic dGVzdDp0ZXN0"})
	assert.Nil(t, err)
}

// clientHTTPBasicAuth validates the basic authentication
type clientHTTPBasicAuth struct{}

func (clientHTTPBasicAuth *clientHTTPBasicAuth) do(req *http.Request) (*http.Response, error) {
	response := http.Response{
		Header: make(http.Header),
		Body:   http.NoBody,
	}
	switch req.Method {
	case "GET":
		if req.Header.Get("Authorization") == "Basic dGVzdDp0ZXN0" {
			response.StatusCode = 200
		} else {
			response.StatusCode = 401
			response.Header.Set("WWW-Authenticate", "Basic realm=\"Haystack\"")
		}
		return &response, nil
	}
	return &response, nil
}
