// Package googp is a OGP (Open Graph protocol) parser library for Golang.
//
// This library is fully compliant with the reference, highly customizable, and supports type conversion.
//
// Links
//
// OGP: https://ogp.me
// Source code: https://github.com/soranoba/googp
//
package googp

import (
	"fmt"
	"mime"
	"net/http"
)

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
			return fmt.Errorf("%w (%s)", ErrUnsupportedPage, mt)
		}
	}

	return NewParser(opts...).Parse(res.Body, i)
}
