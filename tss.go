package tss

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/nathan-fiscaletti/consolesize-go"
	"golang.org/x/net/html"
)

type Template struct {
	rootElement Element
	Elements    map[string]*Element
}

func NewTemplate(tml io.Reader, tss io.Reader) (Template, error) {
	h, err := html.Parse(tml)
	if err != nil {
		panic(err)
	}

	rootElem := h.FirstChild.LastChild.FirstChild
	template := Template{Elements: make(map[string]*Element)}

	var f func(*html.Node) (*Element, error)
	f = func(htmlNode *html.Node) (*Element, error) {
		if htmlNode.Type == html.TextNode && strings.TrimSpace(htmlNode.Data) == "" {
			return nil, nil
		}

		n, err := parseElement(htmlNode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Element: %w", err)
		}

		if id := n.id; id != "" {
			template.Elements[id] = &n
		}

		for c := htmlNode.FirstChild; c != nil; c = c.NextSibling {
			nd, err := f(c)
			if err != nil {
				return nil, err
			}

			if nd != nil {
				n.children = append(n.children, nd)
			}
		}

		return &n, nil
	}
	root, err := f(rootElem)
	if err != nil {
		return Template{}, fmt.Errorf("failed to parse template: %w", err)
	}

	template.rootElement = *root
	return template, nil
}

func (t Template) Render(w io.Writer) {
	cols, _ := consolesize.GetConsoleSize()
	res := t.rootElement.render(cols)
	fmt.Fprint(w, strings.Join(res, "\n")+"\n")
}

func (t Template) RenderFullScreen(w io.Writer) {
	cols, rows := consolesize.GetConsoleSize()
	res := t.rootElement.render(cols)

	b := bufio.NewWriter(w)
	defer b.Flush()
	b.Write([]byte(strings.Join(res, "\n") + strings.Repeat("\n", rows-len(res))))
}
