package googp

import (
	"errors"
	"fmt"
)

var (
	// ErrUnsupportedPage is an unsupported page errror.
	ErrUnsupportedPage = errors.New("Unsupported page")
)

// BadStatusCodeError is an error returned when the status code is not 200 in Fetch.
type BadStatusCodeError struct {
	StatusCode int
}

func (err BadStatusCodeError) Error() string {
	return fmt.Sprintf("Bad status code (%d)", err.StatusCode)
}
