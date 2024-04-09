package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoAuthenticator(t *testing.T) {
	headers, err := NoAuthenticator{}.Authenticate(
		"http://localhost:8080/api/demo/",
		"test",
		"test",
		&clientHTTPNoAuth{},
	)

	assert.Equal(t, headers, map[string]string{})
	assert.Nil(t, err)
}

// clientHTTPNoAuth just returns a 200 for all requests
type clientHTTPNoAuth struct{}

func (clientHTTPNoAuth *clientHTTPNoAuth) do(req *http.Request) (*http.Response, error) {
	response := http.Response{
		Header: make(http.Header),
		Body:   http.NoBody,
	}
	switch req.Method {
	case "GET":
		response.StatusCode = 200
		return &response, nil
	}
	return &response, nil
}
