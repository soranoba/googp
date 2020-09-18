package googp

import (
	"io"
	"reflect"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Meta is a model that structure contents of meta tag in html.
type Meta struct {
	Property string
	Content  string
}

// Parser is an OGP parser.
type Parser struct {
	opts ParserOpts
}

// ParserOpts is an option of `Parser`.
type ParserOpts struct {
	// You can add processing when you need to regard another Nodes as `<meta>`.
	// For example, you can use it when you want to get the `<title>`.
	PreNodeFunc func(*html.Node) *Meta
	// You can add body to parse target.
	// If html have some meta tags in the body, you should set to true.
	IncludeBody bool
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

// ParseNode is execute to parse OGPs from the HTML node.
func (parser *Parser) ParseNode(n *html.Node, i interface{}) error {
	ac := newAccessor(nil, reflect.ValueOf(i))
	return parser.parseNode(n, ac)
}

func (parser *Parser) parseNode(n *html.Node, ac accessor) error {
	switch n.DataAtom {
	case atom.Html, atom.Head, 0:
		return parser.parseChildNode(n, ac)
	case atom.Body:
		if parser.opts.IncludeBody {
			return parser.parseChildNode(n, ac)
		}
	}

	var meta *Meta
	if f := parser.opts.PreNodeFunc; f != nil {
		meta = f(n)
	}
	if meta == nil {
		meta = getOGPMeta(n)
	}

	if meta != nil {
		if err := ac.Set(meta.Property, meta.Content); err != nil {
			return err
		}
	}
	return nil
}

func (parser *Parser) parseChildNode(n *html.Node, ac accessor) error {
	for n := n.FirstChild; n != nil; n = n.NextSibling {
		if err := parser.parseNode(n, ac); err != nil {
			return err
		}
	}
	return nil
}

func getOGPMeta(n *html.Node) *Meta {
	if n.DataAtom != atom.Meta {
		return nil
	}

	meta := new(Meta)
	for _, attr := range n.Attr {
		switch attr.Key {
		case "property":
			meta.Property = attr.Val
		case "content":
			meta.Content = attr.Val
		}
	}

	if meta.Property != "" && meta.Content != "" {
		return meta
	}
	return nil
}
