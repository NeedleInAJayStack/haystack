package client

type AuthError struct {
	message string
}

func NewAuthError(message string) AuthError {
	return AuthError{message: message}
}

func (err AuthError) Error() string {
	return "Auth error: " + err.message
}
