package client

import (
	"github.com/NeedleInAJayStack/haystack"
)

// AuthError represents an error in the authorization
type AuthError struct {
	Message string
}

// NewAuthError creates a new AuthError object.
func NewAuthError(message string) AuthError {
	return AuthError{Message: message}
}

func (err AuthError) Error() string {
	return "Auth error: " + err.Message
}

// CallError occurs when communication is successful, and a Grid is returned, but the grid has an err
// tag indicating a server side error.
type CallError struct {
	Grid haystack.Grid
}

// NewCallError creates a new CallError object.
func NewCallError(grid haystack.Grid) CallError {
	return CallError{Grid: grid}
}

func (err CallError) Error() string {
	dis := err.Grid.Meta().Get("dis")
	switch val := dis.(type) {
	case haystack.Str:
		return "Call error: " + val.String()
	default:
		return "Call error: Server side error"
	}
}

// HTTPError occurs when communication is successful with a server, but we receive an HTTP error response.
type HTTPError struct {
	Code    int
	Message string
}

// NewHTTPError creates a new HTTPError object.
func NewHTTPError(code int, message string) HTTPError {
	return HTTPError{Code: code, Message: message}
}

func (err HTTPError) Error() string {
	return "HTTP error: " + err.Message
}

// NetworkError occurs when there is a network I/O or connection problem with communication to the server.
type NetworkError struct {
	Message string
}

// NewNetworkError creates a new NetworkError object.
func NewNetworkError(message string) NetworkError {
	return NetworkError{Message: message}
}

func (err NetworkError) Error() string {
	return "Network error: " + err.Message
}
