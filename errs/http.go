package errs

import (
	"fmt"
	"io"
	"io/ioutil"
)

type HTTPError struct {
	statusCode int
	message    string
}

func NewHTTPError(code int, body io.ReadCloser) *HTTPError {
	defer body.Close()
	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return &HTTPError{code, ""}
	}
	return &HTTPError{code, string(buf)}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error: %d %s", e.statusCode, e.message)
}
