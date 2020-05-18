package googp

import (
	"reflect"
	"testing"
	"time"
)

func Test_ValueAccessor(t *testing.T) {
	var s string
	field := newValueAccessor(reflect.ValueOf(s)) // string
	assertError(t, field.Set("og:title", "title"))

	field = newValueAccessor(reflect.ValueOf(&s)) // *string
	assertNoError(t, field.Set("og:title", "title"))
	assertEqual(t, s, "title")
	assertNoError(t, field.Set("og:title", "title2"))
	assertEqual(t, s, "title") // cannot overwrite

	var time time.Time
	field = newValueAccessor(reflect.ValueOf(time)) // URL
	assertError(t, field.Set("video:release_date", "2020-05-20T01:01:25Z"))

	field = newValueAccessor(reflect.ValueOf(&time)) // *URL
	assertNoError(t, field.Set("video:release_date", "2020-05-20T01:01:25Z"))
	assertEqual(t, time.String(), "2020-05-20 01:01:25 +0000 UTC")

	var i int
	field = newValueAccessor(reflect.ValueOf(i)) // int
	assertError(t, field.Set("og:number", "123"))

	field = newValueAccessor(reflect.ValueOf(&i)) // *int
	assertNoError(t, field.Set("og:number", "123"))
	assertEqual(t, i, 123)

	var u uint
	field = newValueAccessor(reflect.ValueOf(u)) // uint
	assertError(t, field.Set("og:number", "123"))

	field = newValueAccessor(reflect.ValueOf(&u)) // *uint
	assertNoError(t, field.Set("og:number", "123"))
	assertEqual(t, u, uint(123))

	var f float64
	field = newValueAccessor(reflect.ValueOf(f)) // float64
	assertError(t, field.Set("og:number", "23.5"))

	field = newValueAccessor(reflect.ValueOf(&f)) // *float64
	assertNoError(t, field.Set("og:number", "23.5"))
	assertEqual(t, f, float64(23.5))
}

func Test_ArrayAccessor(t *testing.T) {
	var arr [3]string
	field := newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(arr)) // [3]string
	assertError(t, field.Set("og:image", "http://example.com/image.png"))

	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&arr)) // *[3]string
	assertNoError(t, field.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image2.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image3.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image4.png"))
	assertEqual(t, arr, [3]string{
		"http://example.com/image1.png",
		"http://example.com/image2.png",
		"http://example.com/image3.png",
	})

	var pArr [3]*string
	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(pArr)) // [3]*string
	assertError(t, field.Set("og:image", "http://example.com/image.png"))

	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&pArr)) // *[3]*string
	assertNoError(t, field.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image2.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image3.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image4.png"))
	assertEqual(t, *pArr[0], "http://example.com/image1.png")
	assertEqual(t, *pArr[1], "http://example.com/image2.png")
	assertEqual(t, *pArr[2], "http://example.com/image3.png")

	var sli []string
	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(sli)) // []string
	assertError(t, field.Set("og:image", "http://example.com/image.png"))
	assertEqual(t, len(sli), 0)

	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&sli)) // *[]string
	assertNoError(t, field.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image2.png"))
	assertEqual(t, len(sli), 2)
	assertEqual(t, sli, []string{
		"http://example.com/image1.png",
		"http://example.com/image2.png",
	})

	var pSli []*string
	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(pSli)) // []*string
	assertError(t, field.Set("og:image", "http://example.com/image.png"))

	field = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&pSli)) // *[]*string
	assertNoError(t, field.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, field.Set("og:image", "http://example.com/image2.png"))
	assertEqual(t, len(pSli), 2)
	assertEqual(t, *pSli[0], "http://example.com/image1.png")
	assertEqual(t, *pSli[1], "http://example.com/image2.png")
}

func Test_StructAccessor(t *testing.T) {
	var v struct {
		Title       string  `googp:"og:title"`
		Description *string `googp:"og:description"`
		Url         *string `googp:"og:url"`
	}

	field := newAccessor(nil, reflect.ValueOf(v))
	assertError(t, field.Set("og:title", "title"))
	assertError(t, field.Set("og:description", "description"))

	field = newAccessor(nil, reflect.ValueOf(&v))
	assertNoError(t, field.Set("og:title", "title"))
	assertNoError(t, field.Set("og:description", "description"))
	assertEqual(t, v.Title, "title")
	assertEqual(t, *v.Description, "description")
	assertEqual(t, v.Url, (*string)(nil))
}
