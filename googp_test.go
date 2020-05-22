package googp

import (
	"reflect"
	"runtime"
	"testing"
)

func TestFetch(t *testing.T) {
	ogp := new(OGP)
	assertNoError(t, Fetch(endpoint()+"/1.html", ogp))

	assertEqual(t, ogp.Title, "title")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/image.png")
}

func Test_Fetch_NotFound(t *testing.T) {
	ogp := new(OGP)
	assertError(t, Fetch(endpoint()+"/notfound.html", ogp))
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
