package errs

import "fmt"

type ListedError struct {
	url string
}

func NewListedError(u string) *ListedError {
	return &ListedError{u}
}

func (e *ListedError) Error() string {
	return fmt.Sprintf("%s: %s", "already listed", e.url)
}
