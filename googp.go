package googp

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
)

var (
	BadStatusCodeErr   = errors.New("Bad status code")
	UnsupportedPageErr = errors.New("Unsupported page")
)

// Fetch the content from the URL and parse OGP information.
func Fetch(rawurl string, i interface{}, opts ...ParserOpts) error {
	res, err := http.Get(rawurl)
	if err != nil {
		return fmt.Errorf("Failed to get the content: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%w (%d)", BadStatusCodeErr, res.StatusCode)
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
