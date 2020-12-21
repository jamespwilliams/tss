package tss

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type flow int

const (
	flowRow flow = iota
	flowColumn
)

type width struct {
	value     int
	isPercent bool
}

type Element struct {
	id string

	flow   flow
	width  *width
	border bool

	isContentNode bool
	content       string

	children []*Element
}

func parseElement(node *html.Node) (Element, error) {
	id := getAttributeValue(node, "id")

	flowAttr := getAttributeValue(node, "flow")
	flow := flowColumn
	if flowAttr == "row" {
		flow = flowRow
	}

	border := hasAttribute(node, "border")
	// fmt.Println(border)

	w := 0
	if widthAttr := getAttributeValue(node, "width"); widthAttr != "" {
		var err error
		if w, err = strconv.Atoi(widthAttr); err != nil {
			return Element{}, fmt.Errorf("failed to parse width attribute: %w", err)
		}
	}

	// TODO: make a subtype:
	var content string
	var isContentNode bool
	if node.Type == html.TextNode {
		content = strings.TrimSpace(node.Data)
		isContentNode = true
	}

	return Element{
		id:   id,
		flow: flow,
		width: &width{
			value: w,
		},
		content:       content,
		isContentNode: isContentNode,
		border:        border,
	}, nil
}

func getAttributeValue(node *html.Node, key string) string {
	for _, a := range node.Attr {
		if a.Key == key {
			return a.Val
		}
	}

	return ""
}

func hasAttribute(node *html.Node, key string) bool {
	for _, a := range node.Attr {
		if a.Key == key {
			return true
		}
	}

	return false
}

func (e Element) render(w int) (lines []string) {
	width := e.innerWidth()
	if width == 0 {
		width = w
		if e.border {
			width -= 2
		}
	}

	defer func() {
		if e.border {
			lines = addBorder(lines, width)
		}
		return
	}()

	if e.content != "" {
		// TODO: this would probably be cleaner if this case was a separate type, like textNode or similar
		for start := 0; start < len(e.content); start += width {
			end := start + width
			if end >= len(e.content) {
				end = len(e.content)
			}

			lines = append(lines, e.content[start:end])
		}

		lastLine := lines[len(lines)-1]
		lines[len(lines)-1] = lastLine + strings.Repeat(" ", width-monospaceLength(lastLine))

		return lines
	}

	var childrenLines [][]string
	var longestChildLength int

	for _, child := range e.children {
		lines := child.render(width)
		childrenLines = append(childrenLines, lines)

		if len(lines) > longestChildLength {
			longestChildLength = len(lines)
		}
	}

	if e.flow == flowRow {
		for _, childLines := range childrenLines {
			for _, l := range childLines {
				// fmt.Printf("adding line (row) %v\n", strings.ReplaceAll(l, " ", "~"))

				lines = append(lines, l)
			}
		}
	} else {
		// Iterate for each row, over each child, adding its content to the row:
		for i := 0; i < longestChildLength; i++ {
			var line string
			for childIndex, childLines := range childrenLines {
				child := e.children[childIndex]

				childWidth := child.totalWidth()
				if childWidth == 0 {
					// TODO: what if child is next to other Elements with defined widths?
					childWidth = width
				}

				if i < len(childLines) {
					line += childLines[i]
					if childWidth-monospaceLength(childLines[i]) > 0 {
						line += strings.Repeat(" ", childWidth-monospaceLength(childLines[i]))
					}
				} else {
					line += strings.Repeat(" ", childWidth)
				}
			}

			if monospaceLength(line) < width {
				line += strings.Repeat(" ", width-monospaceLength(line))
			}
			// fmt.Printf("adding line (col) %v\n", strings.ReplaceAll(line, " ", "~"))

			lines = append(lines, line)
		}
	}

	return lines
}

func (c Element) totalWidth() int {
	if c.border {
		return c.width.value + 2
	}

	return c.width.value
}

func (c Element) innerWidth() int {
	return c.width.value
}

func (c *Element) SetContent(content string) {
	c.children = []*Element{&Element{
		width:         &width{},
		content:       content,
		isContentNode: true,
	}}
}

func (c *Element) SetWidth(width int) {
	c.width.value = width
}

func (c Element) String() string {
	return fmt.Sprintf("<width=%v>%v</>", c.width.value, c.content)
}
