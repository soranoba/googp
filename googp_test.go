package googp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	var ogp OGP
	assertNoError(t, Fetch(endpoint()+"/1.html", &ogp))

	assertEqual(t, ogp.Title, "title")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/image.png")
}

func TestFetch_NotFound(t *testing.T) {
	var ogp OGP
	assertError(t, Fetch(endpoint()+"/notfound.html", &ogp))
}

func TestParse(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint()+"/1.html", nil)
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.Do(req.WithContext(ctx))
	var ogp OGP
	assertNoError(t, Parse(res, &ogp))

	assertEqual(t, ogp.Title, "title")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/image.png")
}

func ExampleFetch() {
	var ogp OGP
	if err := Fetch("https://ogp.me", &ogp); err != nil {
		return
	}

	fmt.Printf("og:title = \"%s\"", ogp.Title)
	fmt.Printf("og:type = \"%s\"", ogp.Type)
	fmt.Printf("og:url = \"%s\"", ogp.URL)

	// Outputs:
	// og:title = "Open Graph protocol"
	// og:type = "website"
	// og:url = "https://ogp.me/"
}

func ExampleFetch_customizeModel() {
	type MyOGP struct {
		Title string  `googp:"og:title"`
		URL   url.URL `googp:"og:url"`
		AppID int     `googp:"fb:app_id"`
	}

	var ogp MyOGP
	if err := Fetch("https://ogp.me", &ogp); err != nil {
		return
	}

	fmt.Printf("og:title = \"%s\"", ogp.Title)
	fmt.Printf("og:url = \"%s\"", ogp.URL.String())
	fmt.Printf("fb:app_id = \"%d\"", ogp.AppID)

	// Outputs:
	// og:title = "Open Graph protocol"
	// og:url = "https://ogp.me/"
	// fb:app_id = 115190258555800
}

func assertEqual(t *testing.T, got, expected interface{}) bool {
	if !reflect.DeepEqual(got, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Not equals:\n  file    : %s:%d\n  got     : %#v\n  expected: %#v\n", file, line, got, expected)
		return false
	}
	return true
}

func assertNotEqual(t *testing.T, got, expected interface{}) bool {
	if reflect.DeepEqual(got, expected) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Equals:\n  file    : %s:%d\n  got     : %#v\n", file, line, got)
		return false
	}
	return true
}

func assertError(t *testing.T, err error) bool {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("NoError:\n  file    : %s:%d\n", file, line)
		return false
	}
	return true
}

func assertNoError(t *testing.T, err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Error:\n  file    : %s:%d\n  error   : %#v\n", file, line, err)
		return false
	}
	return true
}
