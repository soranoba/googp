package googp

import "net/http"

// Fetch the content from the URL and parse OGP information.
func Fetch(rawurl string, i interface{}, opts ...ParserOpts) error {
	res, err := http.Get(rawurl)
	if err != nil {
		return err
	}
	return NewParser(opts...).Parse(res.Body, i)
}
