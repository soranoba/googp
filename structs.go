package googp

import (
	"errors"
	"net/url"
)

var (
	ErrInvalidFormat error = errors.New("invalid format")
)

type URL struct {
	url.URL
}

type OGP struct {
	Title  string  `googp:"og:title" json:"title,omitempty"`
	Type   string  `googp:"og:type"  json:"type,omitempty"`
	URL    string  `googp:"og:url"   json:"url,omitempty"`
	Images []Image `googp:"og:image" json:"images,omitempty"`

	Audios          []Audio  `googp:"og:audio"            json:"audios,omitempty"`
	Description     string   `googp:"og:description"      json:"description,omitempty"`
	Determiner      string   `googp:"og:determiner"       json:"determiner,omitempty"`
	Locale          string   `googp:"og:locale"           json:"locale,omitempty"`
	LocaleAlternate []string `googp:"og:locale:alternate" json:"locale_alternate,omitempty"`
	SiteName        string   `googp:"og:site_name"        json:"site_name,omitempty"`
	Videos          []Video  `googp:"og:video"            json:"videos,omitempty"`
}

type Image struct {
	URL       string `googp:"og:image,og:image:url" json:"url,omitempty"`
	SecureURL string `googp:"og:image:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:image:type"         json:"type,omitempty"`
	Width     int    `googp:"og:image:width"        json:"width,omitempty"`
	Height    int    `googp:"og:image:height"       json:"height,omitempty"`
	Alt       string `googp:"og:image:alt"          json:"alt,omitempty"`
}

type Audio struct {
	URL       string `googp:"og:audio,og:audio:url" json:"url,omitempty"`
	SecureURL string `googp:"og:audio:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:audio:type"         json:"type,omitempty"`
}

type Video struct {
	URL       string `googp:"og:video,og:video:url" json:"url,omitempty"`
	SecureURL string `googp:"og:video:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:video:type"         json:"type,omitempty"`
	Width     int    `googp:"og:video:width"        json:"width,omitempty"`
	Height    int    `googp:"og:video:height"       json:"height,omitempty"`
}

func (url URL) MarshalText() (text []byte, err error) {
	return []byte(url.String()), nil
}

func (url *URL) UnmarshalText(text []byte) error {
	u, err := url.Parse(string(text))
	if err != nil {
		return err
	}
	url.URL = *u
	return nil
}

func (url URL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + url.String() + "\""), nil
}

func (u *URL) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return ErrInvalidFormat
	}

	if !(data[0] == '"' && data[len(data)-1] == '"') {
		return ErrInvalidFormat
	}

	got, err := url.Parse(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	u.URL = *got
	return nil
}
