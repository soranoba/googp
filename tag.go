package googp

import (
	"reflect"
	"strings"
)

const (
	structTagKey = "googp"
)

// tag is struct tag supported by googp
type tag struct {
	// Array of OGP property names. (e.g. `og:title`)
	names []string
}

// newTag is create a `*tag` from `reflect.StructField`
func newTag(f reflect.StructField) *tag {
	value := f.Tag.Get(structTagKey)
	var names []string

	if value == "" {
		// NOTE: If tag is not specified, it is same as being given `og:${field_name}`.
		names = []string{"og:" + toSnake(f.Name)}
	} else {
		names = strings.Split(value, ",")
	}
	return &tag{names: names}
}

// isContainsName returns true, when the tag contains the name.
// Otherwise, it returns false.
func (t *tag) isContainsName(name string) bool {
	for _, n := range t.names {
		if n == name {
			return true
		}
	}
	return false
}
