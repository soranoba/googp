package googp

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
)

var (
	UnsupportedPageErr = errors.New("Unsupported page")
)

// BadStatusCodeError is an error returned when the status code is not 200 in Fetch.
type BadStatusCodeError struct {
	StatusCode int
}

func (err BadStatusCodeError) Error() string {
	return fmt.Sprintf("Bad status code (%d)", err.StatusCode)
}

// Fetch the content from the URL and parse OGP information.
func Fetch(rawurl string, i interface{}, opts ...ParserOpts) error {
	res, err := http.Get(rawurl)
	if err != nil {
		return fmt.Errorf("Failed to get the content: %w", err)
	}
	return Parse(res, i, opts...)
}

// Parse OGP information.
// It returns an error when the status code of the response is error.
func Parse(res *http.Response, i interface{}, opts ...ParserOpts) error {
	if res.StatusCode != 200 {
		return &BadStatusCodeError{StatusCode: res.StatusCode}
	}

	ct := res.Header.Get("Content-Type")
	if ct != "" {
		mt, _, err := mime.ParseMediaType(ct)
		if err != nil {
			return fmt.Errorf("Invalid Content-Type: %w", err)
		}
		if mt != "text/html" {
			return fmt.Errorf("%w (%s)", UnsupportedPageErr, mt)
		}
	}

	return NewParser(opts...).Parse(res.Body, i)
}
