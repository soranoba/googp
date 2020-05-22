# googp
[![CircleCI](https://circleci.com/gh/soranoba/googp.svg?style=svg&circle-token=4282311988e7373cbc6033771566d912f1f446c9)](https://circleci.com/gh/soranoba/googp)
[![Go Report Card](https://goreportcard.com/badge/github.com/soranoba/googp)](https://goreportcard.com/report/github.com/soranoba/googp)
[![GoDoc](https://godoc.org/github.com/soranoba/googp?status.svg)](https://godoc.org/github.com/soranoba/googp)

googp is a [OGP (Open Graph protocol)](https://ogp.me/) parser library for Golang.

## Overviews

- Mapping to your struct and OGP tags.
- Supports all types that implement [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler).

## Usage

```go
import (
    "fmt"
    "github.com/soranoba/googp"
)

type CustomOGP struct {
    Title       string   `googp:"og:title"`
    Description string   `googp:"-"`        // ignored
    images      []string                    // private field ignored
    Videos      []string `googp:"og:video,og:video:url"`
    Musics      Music    `googp:"music"`    // object type
}

func main() {
    ogp1 := new(googp.OGP)
    if err := googp.Fetch("https://soranoba.net", ogp1); err != nil {
        return
    }
    fmt.Println(ogp1)

    ogp2 := new(CustomOGP)
    if err := googp.Fetch("https://soranoba.net", ogp2); err != nil {
        return
    }
    fmt.Println(ogp2)
}
```

## Object Mappings

### [Structured Properties](https://ogp.me/#structured)

```go
type OGP struct {
    Image struct {
        URL       `googp:"og:image,og:image:url"`
        SecureURL `googp:"og:image:secure_url"`
    } `googp:"og:image"`
}
```

You may collect in a struct by specifying the root tag.<br>
In case of specifying `og:image`, googp collect values which property is `og:image:*`.

### [Arrays](https://ogp.me/#array)

```go
type OGP struct {
    Image []string `googp:"og:image"`
}
```

googp collects values which the same properties.

### [Object Types](https://ogp.me/#types)

In googp, it same as Structured Properties.<br>
You may define your own type yourself.