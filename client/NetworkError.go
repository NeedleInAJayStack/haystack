package client

type NetworkError struct {
	message string
}

func NewNetworkError(message string) NetworkError {
	return NetworkError{message: message}
}

func (err NetworkError) Error() string {
	return err.message
}
