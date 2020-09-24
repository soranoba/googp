package googp

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestParser_Parse_1(t *testing.T) {
	res, err := http.Get(endpoint() + "/1.html")
	assertNoError(t, err)

	parser := NewParser()
	var ogp OGP
	assertNoError(t, parser.Parse(res.Body, &ogp))

	assertEqual(t, ogp.Title, "title")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/image.png")
}

func TestParser_Parse_1_PreNodeFunc(t *testing.T) {
	res, err := http.Get(endpoint() + "/1.html")
	assertNoError(t, err)

	parser := NewParser(ParserOpts{PreNodeFunc: func(node *html.Node) *Meta {
		if node.DataAtom == atom.Title {
			return &Meta{Property: "og:title", Content: node.FirstChild.Data}
		}
		return nil
	}})
	var ogp OGP
	assertNoError(t, parser.Parse(res.Body, &ogp))

	assertEqual(t, ogp.Title, "SamplePage")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/image.png")
}

func TestParser_Parse_2(t *testing.T) {
	res, err := http.Get(endpoint() + "/2.html")
	assertNoError(t, err)

	parser := NewParser()
	var ogp OGP
	assertNoError(t, parser.Parse(res.Body, &ogp))

	assertEqual(t, ogp.Title, "title")
	assertEqual(t, ogp.Type, "website")
	assertEqual(t, ogp.URL, "http://example.com")
	assertEqual(t, ogp.Images[0].URL, "http://example.com/rock.jpg")
	assertEqual(t, ogp.Images[0].SecureURL, "https://secure.example.com/ogp.jpg")
	assertEqual(t, ogp.Images[0].Type, "image/jpeg")
	assertEqual(t, ogp.Images[0].Width, 400)
	assertEqual(t, ogp.Images[0].Height, 300)
	assertEqual(t, ogp.Images[0].Alt, "A shiny red apple with a bite taken out")
	assertEqual(t, ogp.Images[1].URL, "http://example.com/rock2.jpg")
	assertEqual(t, ogp.Images[2].URL, "http://example.com/rock3.jpg")
	assertEqual(t, ogp.Images[2].Height, 1000)
	assertEqual(t, ogp.Audios[0].URL, "http://example.com/sound.mp3")
	assertEqual(t, ogp.Audios[0].SecureURL, "https://secure.example.com/sound.mp3")
	assertEqual(t, ogp.Audios[0].Type, "audio/mpeg")
	assertEqual(t, ogp.Videos[0].URL, "http://example.com/movie.swf")
	assertEqual(t, ogp.Videos[0].SecureURL, "https://secure.example.com/movie.swf")
	assertEqual(t, ogp.Videos[0].Type, "application/x-shockwave-flash")
	assertEqual(t, ogp.Videos[0].Width, 400)
	assertEqual(t, ogp.Videos[0].Height, 300)
	assertEqual(t, ogp.Audios[1].URL, "http://example.com/bond/theme.mp3")
	assertEqual(t, ogp.Description, "Sean Connery found fame and fortune as the suave, sophisticated British agent, James Bond.")
	assertEqual(t, ogp.Determiner, "the")
	assertEqual(t, ogp.Locale, "en_GB")
	assertEqual(t, ogp.LocaleAlternate[0], "fr_FR")
	assertEqual(t, ogp.LocaleAlternate[1], "es_ES")
	assertEqual(t, ogp.SiteName, "IMDb")
}

func TestParser_Parse_2_Custom(t *testing.T) {
	res, err := http.Get(endpoint() + "/2.html")
	assertNoError(t, err)

	parser := NewParser()
	type CustomOGP struct {
		Images []string `googp:"og:image"`
	}
	var ogp CustomOGP
	assertNoError(t, parser.Parse(res.Body, &ogp))

	assertEqual(t, ogp.Images[0], "http://example.com/rock.jpg")
	assertEqual(t, ogp.Images[1], "http://example.com/rock2.jpg")
	assertEqual(t, ogp.Images[2], "http://example.com/rock3.jpg")
}

func TestParser_Parse_3(t *testing.T) {
	res, err := http.Get(endpoint() + "/3.html")
	assertNoError(t, err)

	parser := NewParser()
	var ogp OGP
	assertEqual(
		t,
		fmt.Sprintf("%+v", parser.Parse(res.Body, &ogp)),
		"og:image:width field is invalid. (type = int, value = invalid)",
	)
}

func TestParser_Parse_4(t *testing.T) {
	res, err := http.Get(endpoint() + "/4.html")
	assertNoError(t, err)

	parser := NewParser(ParserOpts{IncludeBody: true})
	var ogp OGP
	assertNoError(t, parser.Parse(res.Body, &ogp))

	assertEqual(t, ogp.Title, "og title")
	assertEqual(t, ogp.Type, "video.other")
	assertEqual(t, ogp.URL, "https://example.com/url")
	assertEqual(t, ogp.Images[0].URL, "https://example.com/image")
}

func ExampleParser_Parse() {
	reader := strings.NewReader(`
		<html>
			<head>
				<meta property="og:title" content="title" />
				<meta property="og:type" content="website" />
				<meta property="og:url" content="http://example.com" />
				<meta property="og:image" content="http://example.com/image.png" />
			</head>
			<body>
			</body>
		</html>
	`)

	var ogp OGP
	if err := NewParser().Parse(reader, &ogp); err != nil {
		return
	}

	fmt.Printf("og:title = \"%s\"\n", ogp.Title)
	fmt.Printf("og:type = \"%s\"\n", ogp.Type)
	fmt.Printf("og:url = \"%s\"\n", ogp.URL)
	fmt.Printf("og:image = \"%s\"\n", ogp.Images[0].URL)

	// Output:
	// og:title = "title"
	// og:type = "website"
	// og:url = "http://example.com"
	// og:image = "http://example.com/image.png"
}

func endpoint() string {
	str, ok := os.LookupEnv("NGINX_HOST")
	if ok {
		return "http://" + str
	}
	return "http://localhost:8080"
}
