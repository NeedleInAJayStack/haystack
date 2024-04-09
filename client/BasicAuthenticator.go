package client

import (
	"encoding/base64"
	"net/http"
)

// BasicAuthenticator is an Authenticator that uses basic authentication
type BasicAuthenticator struct{}

func (authenticator BasicAuthenticator) Authenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
) (map[string]string, error) {
	authValue := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	basicAuth := "Basic " + authValue
	headers := map[string]string{"Authorization": basicAuth}

	// Test the basic auth to ensure that it works
	req, _ := http.NewRequest("GET", uri, nil)
	setHeaders(req, headers)
	resp, _ := client.do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return map[string]string{}, NewAuthError("Basic auth failed with status: " + resp.Status)
	}

	return headers, nil
}
