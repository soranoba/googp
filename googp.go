// Package googp is a OGP (Open Graph protocol) parser library for Golang.
//
// This library is fully compliant with the reference, highly customizable, and supports type conversion.
//
// Links
//
// OGP: https://ogp.me
//
// Source code: https://github.com/soranoba/googp
//
package googp

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"net/http"

	"golang.org/x/net/html/charset"
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

	br := bufio.NewReader(res.Body)
	var reader io.Reader = br
	data, _ := br.Peek(1024)
	enc, _, _ := charset.DetermineEncoding(data, ct)
	reader = enc.NewDecoder().Reader(reader)

	return NewParser(opts...).Parse(reader, i)
}
