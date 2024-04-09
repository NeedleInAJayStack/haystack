package client

// NoAuthenticator is an Authenticator that simply assumes the client is already authenticated
type NoAuthenticator struct{}

func (authenticator NoAuthenticator) Authenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
) (map[string]string, error) {
	return map[string]string{}, nil
}
