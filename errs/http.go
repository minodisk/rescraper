package errs

import (
	"fmt"
)

type HTTPError struct {
	statusCode int
	message    string
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{code, message}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error: %d %s", e.statusCode, e.message)
}
