package googp

// OGP is a model that have Basic Metadata and Optional Metadata defined in the reference.
// ref: https://ogp.me/
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

// Image is a model that structure contents of og:image.
type Image struct {
	URL       string `googp:"og:image,og:image:url" json:"url,omitempty"`
	SecureURL string `googp:"og:image:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:image:type"         json:"type,omitempty"`
	Width     int    `googp:"og:image:width"        json:"width,omitempty"`
	Height    int    `googp:"og:image:height"       json:"height,omitempty"`
	Alt       string `googp:"og:image:alt"          json:"alt,omitempty"`
}

// Audio is a model that structure contents of og:audio.
type Audio struct {
	URL       string `googp:"og:audio,og:audio:url" json:"url,omitempty"`
	SecureURL string `googp:"og:audio:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:audio:type"         json:"type,omitempty"`
}

// Video is a model that structure contents of og:video.
type Video struct {
	URL       string `googp:"og:video,og:video:url" json:"url,omitempty"`
	SecureURL string `googp:"og:video:secure_url"   json:"secure_url,omitempty"`
	Type      string `googp:"og:video:type"         json:"type,omitempty"`
	Width     int    `googp:"og:video:width"        json:"width,omitempty"`
	Height    int    `googp:"og:video:height"       json:"height,omitempty"`
}
