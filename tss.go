package tss

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nathan-fiscaletti/consolesize-go"
	"golang.org/x/net/html"
)

type Template struct {
	n node
}

func NewTemplate(tml io.Reader, tss io.Reader) (Template, error) {
	h, err := html.Parse(tml)
	if err != nil {
		panic(err)
	}

	rootElem := h.FirstChild.LastChild.FirstChild

	var f func(*html.Node) *node
	f = func(htmlNode *html.Node) *node {
		if htmlNode.Type == html.TextNode && strings.TrimSpace(htmlNode.Data) == "" {
			return nil
		}

		n := convert(htmlNode)

		for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
			nd := f(c)
			if nd != nil {
				n.children = append(n.children, *nd)
			}
		}

		return &n
	}
	rootNode := *f(rootElem)

	return Template{
		n: rootNode,
	}, nil
}

func (t Template) Render(w io.Writer) {
	cols, _ := consolesize.GetConsoleSize()
	res := t.n.render(cols)
	fmt.Fprint(w, strings.Join(res, "\n"))
}

func convert(n *html.Node) node {
	widthAttr := getAttributeValue(n, "width")
	flowAttr := getAttributeValue(n, "flow")

	// TODO(jpw): check err
	w, _ := strconv.Atoi(widthAttr)

	flow := flowColumn
	if flowAttr == "row" {
		flow = flowRow
	}

	var content string
	if n.Type == html.TextNode {
		content = strings.TrimSpace(n.Data)
	}

	return node{
		f: flow,
		w: width{
			value: w,
		},
		content: content,
	}
}

func getAttributeValue(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}
