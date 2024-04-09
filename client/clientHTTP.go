package client

import (
	"net/http"
)

// ClientHTTP wraps http.client in an interface to allow dependency-injection testing
type ClientHTTP interface {
	// Perform an HTTP request and return the response
	do(req *http.Request) (*http.Response, error)
}

// clientHTTPImpl is the default implementation of clientHTTP
type clientHTTPImpl struct {
	httpClient *http.Client
}

func (clientHTTP *clientHTTPImpl) do(req *http.Request) (*http.Response, error) {
	return clientHTTP.httpClient.Do(req)
}
