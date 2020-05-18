package googp

import (
	"reflect"
	"testing"
)

func Test_Tag(t *testing.T) {
	var v1 struct {
		A string `googp:"og:title,og:description"`
	}
	tag := newTag(reflect.TypeOf(v1).Field(0))
	assertEqual(t, tag.names, []string{"og:title", "og:description"})
	assertEqual(t, tag.isContainsName("og:title"), true)
	assertEqual(t, tag.isContainsName("og:description"), true)
	assertEqual(t, tag.isContainsName("og"), false)

	var v2 struct {
		A string
	}
	tag = newTag(reflect.TypeOf(v2).Field(0))
	assertEqual(t, tag.names, []string{"og:a"})
}
