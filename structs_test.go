package googp

import (
	"encoding"
	"encoding/json"
	"testing"
)

func TestURL_JSONMarshaler(t *testing.T) {
	var model struct {
		URL URL `json:"url"`
	}

	_ = json.Marshaler(&model.URL)
	_ = json.Unmarshaler(&model.URL)

	assertNoError(t, json.Unmarshal([]byte("{\"url\":\"https://example.com/hoge/fuga\"}"), &model))
	assertEqual(t, model.URL.String(), "https://example.com/hoge/fuga")

	b, err := json.Marshal(model)
	assertNoError(t, err)
	assertEqual(t, string(b), "{\"url\":\"https://example.com/hoge/fuga\"}")
}

func TestURL_TextMarshaler(t *testing.T) {
	var model struct {
		URL URL
	}

	_ = encoding.TextMarshaler(&model.URL)
	_ = encoding.TextUnmarshaler(&model.URL)
	assertNoError(t, model.URL.UnmarshalText([]byte("https://example.com/hoge/fuga")))

	b, err := model.URL.MarshalText()
	assertNoError(t, err)
	assertEqual(t, string(b), "https://example.com/hoge/fuga")
}
