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

	switch value {
	case "-":
		names = []string{}
	case "":
		if f.Anonymous {
			names = []string{""}
		} else {
			// NOTE: If tag is not specified, it is same as being given `og:${field_name}`.
			names = []string{"og:" + toSnake(f.Name)}
		}
	default:
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

// toSnake converts the string into a snake case.
func toSnake(str string) string {
	runes := []rune(str)
	var p int
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if c >= 'A' && c <= 'Z' {
			runes[i] = c - ('A' - 'a')
			if p+1 < i {
				tmp := append([]rune{'_'}, runes[i:]...)
				runes = append(runes[0:i], tmp...)
			}
			p = i
		}
	}
	return string(runes)
}
