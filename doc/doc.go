package doc

import (
	"bytes"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Doc struct {
	*html.Node
}

func (d *Doc) Query(query string) *Doc {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return &Doc{}
	}
	node := cascadia.Query(d.Node, sel)
	if node == nil {
		return &Doc{Node: &html.Node{}}
	}
	return &Doc{Node: node}
}

func (d *Doc) QueryAll(query string) []*Doc {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return []*Doc{}
	}
	docs := []*Doc{}
	for _, node := range cascadia.QueryAll(d.Node, sel) {
		docs = append(docs, &Doc{Node: node})
	}
	return docs
}

func (d *Doc) NthChild(n int) *Doc {
	p := 0
	for child := d.Node.FirstChild; child != nil; child = child.NextSibling {
		if p == n {
			return &Doc{Node: child}
		}
		p++
	}
	return &Doc{}
}

func (d *Doc) Text() string {
	var buf bytes.Buffer
	var f func(node *html.Node)
	f = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(d.Node)
	return buf.String()
}

func (d *Doc) AttrOr(attr, or string) string {
	for _, a := range d.Node.Attr {
		if a.Key == attr {
			return a.Val
		}
	}
	return or
}
