package googp

import (
	"reflect"
	"testing"
)

func Test_Tag(t *testing.T) {
	var v struct {
		A string `googp:"og:title,og:description"`
		B string
		C string `googp:"-"`
		OGP
	}

	tag := newTag(reflect.TypeOf(v).Field(0))
	assertEqual(t, tag.names, []string{"og:title", "og:description"})
	assertEqual(t, tag.isContainsName("og:title"), true)
	assertEqual(t, tag.isContainsName("og:description"), true)
	assertEqual(t, tag.isContainsName("og"), false)

	tag = newTag(reflect.TypeOf(v).Field(1))
	assertEqual(t, tag.names, []string{"og:b"})

	tag = newTag(reflect.TypeOf(v).Field(2))
	assertEqual(t, tag.names, []string{})

	tag = newTag(reflect.TypeOf(v).Field((3)))
	assertEqual(t, tag.names, []string{""})
}
