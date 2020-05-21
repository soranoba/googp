# googp
[![CircleCI](https://circleci.com/gh/soranoba/googp.svg?style=svg&circle-token=4282311988e7373cbc6033771566d912f1f446c9)](https://circleci.com/gh/soranoba/googp)
[![Go Report Card](https://goreportcard.com/badge/github.com/soranoba/googp)](https://goreportcard.com/report/github.com/soranoba/googp)
[![GoDoc](https://godoc.org/github.com/soranoba/googp?status.svg)](https://godoc.org/github.com/soranoba/googp)

googp is a OGP (Open Graph protocol) parser library for Golang.

## Overviews

- Mapping to your struct and OGP tags.
- Supports all types that implement [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler).

## Usage

```go
import (
    "github.com/soranoba/googp"
)

type CustomOGP struct {
    Title string `googp:"og:title"`
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
