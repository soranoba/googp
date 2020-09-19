package googp

import (
	"reflect"
	"testing"
	"time"
)

func Test_ValueAccessor(t *testing.T) {
	var s string
	ac := newAccessor(nil, reflect.ValueOf(s)) // string
	assertError(t, ac.Set("og:title", "title"))

	ac = newAccessor(nil, reflect.ValueOf(&s)) // *string
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, s, "title")
	assertNoError(t, ac.Set("og:title", "title2"))
	assertEqual(t, s, "title") // cannot overwrite

	var time time.Time
	ac = newAccessor(nil, reflect.ValueOf(time)) // Time
	assertError(t, ac.Set("video:release_date", "2020-05-20T01:01:25Z"))

	ac = newAccessor(nil, reflect.ValueOf(&time)) // *Time
	assertNoError(t, ac.Set("video:release_date", "2020-05-20T01:01:25Z"))
	assertEqual(t, time.String(), "2020-05-20 01:01:25 +0000 UTC")

	var i int
	ac = newValueAccessor(reflect.ValueOf(i)) // int
	assertError(t, ac.Set("og:number", "123"))

	ac = newValueAccessor(reflect.ValueOf(&i)) // *int
	assertNoError(t, ac.Set("og:number", "123"))
	assertEqual(t, i, 123)

	var u uint
	ac = newValueAccessor(reflect.ValueOf(u)) // uint
	assertError(t, ac.Set("og:number", "123"))

	ac = newValueAccessor(reflect.ValueOf(&u)) // *uint
	assertNoError(t, ac.Set("og:number", "123"))
	assertEqual(t, u, uint(123))

	var f float64
	ac = newAccessor(nil, reflect.ValueOf(f)) // float64
	assertError(t, ac.Set("og:number", "23.5"))

	ac = newAccessor(nil, reflect.ValueOf(&f)) // *float64
	assertNoError(t, ac.Set("og:number", "23.5"))
	assertEqual(t, f, float64(23.5))
}

func Test_ArrayAccessor(t *testing.T) {
	var arr [3]string
	ac := newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(arr)) // [3]string
	assertError(t, ac.Set("og:image", "http://example.com/image.png"))

	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&arr)) // *[3]string
	assertNoError(t, ac.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image2.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image3.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image4.png"))
	assertEqual(t, arr, [3]string{
		"http://example.com/image1.png",
		"http://example.com/image2.png",
		"http://example.com/image3.png",
	})

	var pArr [3]*string
	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(pArr)) // [3]*string
	assertError(t, ac.Set("og:image", "http://example.com/image.png"))

	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&pArr)) // *[3]*string
	assertNoError(t, ac.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image2.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image3.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image4.png"))
	assertEqual(t, *pArr[0], "http://example.com/image1.png")
	assertEqual(t, *pArr[1], "http://example.com/image2.png")
	assertEqual(t, *pArr[2], "http://example.com/image3.png")

	var sli []string
	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(sli)) // []string
	assertError(t, ac.Set("og:image", "http://example.com/image.png"))
	assertEqual(t, len(sli), 0)

	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&sli)) // *[]string
	assertNoError(t, ac.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image2.png"))
	assertEqual(t, len(sli), 2)
	assertEqual(t, sli, []string{
		"http://example.com/image1.png",
		"http://example.com/image2.png",
	})

	var pSli []*string
	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(pSli)) // []*string
	assertError(t, ac.Set("og:image", "http://example.com/image.png"))

	ac = newAccessor(&tag{names: []string{"og:image"}}, reflect.ValueOf(&pSli)) // *[]*string
	assertNoError(t, ac.Set("og:image", "http://example.com/image1.png"))
	assertNoError(t, ac.Set("og:image", "http://example.com/image2.png"))
	assertEqual(t, len(pSli), 2)
	assertEqual(t, *pSli[0], "http://example.com/image1.png")
	assertEqual(t, *pSli[1], "http://example.com/image2.png")
}

func Test_StructAccessor(t *testing.T) {
	var v struct {
		Title       string  `googp:"og:title"`
		Description *string `googp:"og:description"`
		URL         *string `googp:"og:url"`
	}

	ac := newAccessor(nil, reflect.ValueOf(v))
	assertError(t, ac.Set("og:title", "title"))
	assertError(t, ac.Set("og:description", "description"))

	ac = newAccessor(nil, reflect.ValueOf(&v))
	assertNoError(t, ac.Set("og:title", "title"))
	assertNoError(t, ac.Set("og:description", "description"))
	assertEqual(t, v.Title, "title")
	assertEqual(t, *v.Description, "description")
	assertEqual(t, v.URL, (*string)(nil))
}

func Test_StructAccessor_ConflictTag(t *testing.T) {
	var v struct {
		A string `googp:"og:title"`
		B string `googp:"og:title"`
	}

	ac := newAccessor(nil, reflect.ValueOf(&v))
	assertNoError(t, ac.Set("og:title", "title"))
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, v.A, "title")
	// If there is a conflict, the backward will be ignored.
	assertEqual(t, v.B, "")
}

func Test_StructAccessor_NoTag(t *testing.T) {
	var v struct {
		A string
	}

	ac := newAccessor(nil, reflect.ValueOf(&v))
	assertNoError(t, ac.Set("og:a", "title"))
	// default is treated as `og:${acName}`
	assertEqual(t, v.A, "title")
}

func Test_StructAccessor_MultipleNames(t *testing.T) {
	var v struct {
		A string `googp:"og:title,og:description"`
	}

	ac := newAccessor(nil, reflect.ValueOf(&v))
	assertNoError(t, ac.Set("og:title", "title"))
	assertNoError(t, ac.Set("og:description", "description"))
	// The one set first has priority
	assertEqual(t, v.A, "title")
}

func Test_StructAccessor_Anonymous(t *testing.T) {
	var og1 struct {
		OGP
	}
	ac := newAccessor(nil, reflect.ValueOf(&og1))
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, og1.Title, "title")

	var og2 struct {
		*OGP
	}
	ac = newAccessor(nil, reflect.ValueOf(&og2))
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, og2.Title, "title")
}

func Test_StructAccessor_Private(t *testing.T) {
	var og1 struct {
		ogp OGP
		A   string
	}
	ac := newAccessor(nil, reflect.ValueOf(&og1))
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, og1.ogp.Title, "")

	var og2 struct {
		ogp *OGP
		A   string
	}
	ac = newAccessor(nil, reflect.ValueOf(&og2))
	assertNoError(t, ac.Set("og:title", "title"))
	assertEqual(t, og2.ogp, (*OGP)(nil))
}
