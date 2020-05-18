package googp

import (
	"io"
	"reflect"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Parser is an OGP parser.
type Parser struct {
	opts ParserOpts
}

// ParserOpts is option of `Parser`.
type ParserOpts struct {
}

type meta struct {
	Property string
	Content  string
}

// NewParser create a `Parser`
func NewParser(opts ...ParserOpts) *Parser {
	switch len(opts) {
	case 0:
		return &Parser{opts: ParserOpts{}}
	case 1:
		return &Parser{opts: opts[0]}
	default:
		panic("Cannot specify multiple ParserOpts")
	}
}

// Parse OGPs from the HTML.
func (parser *Parser) Parse(reader io.Reader, i interface{}) error {
	node, err := html.Parse(reader)
	if err != nil {
		return err
	}
	return parser.ParseNode(node, i)
}

// Parse OGPs from the HTML node.
func (parser *Parser) ParseNode(n *html.Node, i interface{}) error {
	ac := newAccessor(nil, reflect.ValueOf(i))
	return parser.parseNode(n, ac)
}

func (parser *Parser) parseNode(n *html.Node, ac accessor) error {
	for n := n.FirstChild; n != nil; n = n.NextSibling {
		if n.DataAtom == atom.Meta {
			if meta := getOGPMeta(n.Attr); meta != nil {
				if err := ac.Set(meta.Property, meta.Content); err != nil {
					return err
				}
			}
		} else {
			if err := parser.parseNode(n, ac); err != nil {
				return err
			}
		}
	}
	return nil
}

func getOGPMeta(attr []html.Attribute) *meta {
	meta := new(meta)
	for _, a := range attr {
		switch a.Key {
		case "property":
			meta.Property = a.Val
		case "content":
			meta.Content = a.Val
		}
	}

	if meta.Property != "" && meta.Content != "" {
		return meta
	}
	return nil
}
