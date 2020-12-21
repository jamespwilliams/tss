package tss

import (
	"fmt"
	"io"
	"strings"

	"github.com/nathan-fiscaletti/consolesize-go"
	"golang.org/x/net/html"
)

type Template struct {
	rootElement element
}

func NewTemplate(tml io.Reader, tss io.Reader) (Template, error) {
	h, err := html.Parse(tml)
	if err != nil {
		panic(err)
	}

	rootElem := h.FirstChild.LastChild.FirstChild

	var f func(*html.Node) (*element, error)
	f = func(htmlNode *html.Node) (*element, error) {
		if htmlNode.Type == html.TextNode && strings.TrimSpace(htmlNode.Data) == "" {
			return nil, nil
		}

		n, err := parseElement(htmlNode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse element: %w", err)
		}

		for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
			nd, err := f(c)
			if err != nil {
				return nil, err
			}

			if nd != nil {
				n.children = append(n.children, *nd)
			}
		}

		return &n, nil
	}
	root, err := f(rootElem)
	if err != nil {
		return Template{}, fmt.Errorf("failed to parse template: %w", err)
	}

	return Template{
		rootElement: *root,
	}, nil
}

func (t Template) Render(w io.Writer) {
	cols, _ := consolesize.GetConsoleSize()
	res := t.rootElement.render(cols)
	fmt.Fprint(w, strings.Join(res, "\n")+"\n")
}
