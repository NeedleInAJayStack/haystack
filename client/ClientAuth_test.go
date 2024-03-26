package client

import (
	"net/http"
	"testing"
)

func TestClientAuth_NoAuth(t *testing.T) {
	client := &Client{
		clientHTTP: &clientHTTPNoAuth{},
		uri:        "http://localhost:8080/api/demo/",
		username:   "test",
		password:   "test",
	}
	openErr := client.Open()
	if openErr != nil {
		t.Error(openErr)
	}
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

func TestClientAuth_BasicAuth(t *testing.T) {
	client := &Client{
		clientHTTP: &clientHTTPBasicAuth{},
		uri:        "http://localhost:8080/api/demo/",
		username:   "test",
		password:   "test",
	}
	openErr := client.Open()
	if openErr != nil {
		t.Error(openErr)
	}
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
