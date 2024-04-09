package client

// Authenticator is an interface for authenticating a client
type Authenticator interface {
	// Authenticates the client and returns the headers to be used in subsequent requests.
	Authenticate(
		uri string,
		username string,
		password string,
		client ClientHTTP,
	) (map[string]string, error)
}
