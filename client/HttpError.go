package client

type HttpError struct {
	code    int
	message string
}

func NewHttpError(code int, message string) HttpError {
	return HttpError{code: code, message: message}
}

func (err HttpError) Error() string {
	return "HTTP error: " + string(err.code) + err.message
}
