package errs

import "strings"

type MultiError []error

func NewMultiError() MultiError {
	return MultiError{}
}

func (e MultiError) Append(err error) MultiError {
	return append(e, err)
}

func (e MultiError) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "\n")
}
