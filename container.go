package tss

import (
	"fmt"
	"strings"
	"unicode/utf8"
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

type node struct {
	f flow
	w width

	content string

	children []node
}

func (n node) render(w int) (lines []string) {
	width := n.w.value
	if width == 0 {
		width = w
	}

	width -= 2
	defer func() {
		var res []string
		res = append(res, "┌"+strings.Repeat("─", width)+"┐")
		fmt.Println("added            ", res[len(res)-1])
		fmt.Println(width)
		for _, line := range lines {
			res = append(res, "│"+line+"│")
		}
		res = append(res, "└"+strings.Repeat("─", width)+"┘")
		lines = res
		return
	}()

	if n.content != "" {
		// TODO: this would probably be cleaner if this case was a separate type, like textNode or similar
		for start := 0; start < len(n.content); start += width {
			end := start + width
			if end >= len(n.content) {
				end = len(n.content)
			}

			lines = append(lines, n.content[start:end])
		}

		lastLine := lines[len(lines)-1]
		lines[len(lines)-1] = lastLine + strings.Repeat(" ", width-monospaceLength(lastLine))

		return lines
	}

	var childrenLines [][]string
	var longestChildLength int

	for _, child := range n.children {
		lines := child.render(width)
		childrenLines = append(childrenLines, lines)

		if len(lines) > longestChildLength {
			longestChildLength = len(lines)
		}
	}

	// TODO(jpw): fix all these terribly short variable names...
	if n.f == flowRow {
		for _, childLines := range childrenLines {
			for _, l := range childLines {
				fmt.Printf("adding line (row) %v\n", strings.ReplaceAll(l, " ", "~"))

				lines = append(lines, l)
			}
		}
	} else {
		// Iterate for each row, over each child, adding its content to the row:
		for i := 0; i < longestChildLength; i++ {
			var line string
			for childIndex, childLines := range childrenLines {
				child := n.children[childIndex]

				childWidth := child.w.value
				if childWidth == 0 {
					// TODO: what if child is next to other elements with defined widths?
					childWidth = width
				}

				if i < len(childLines) {
					line += childLines[i]
					if childWidth-monospaceLength(childLines[i]) > 0 {
						line += strings.Repeat(" ", childWidth-monospaceLength(childLines[i]))
					}
				} else {
					line += strings.Repeat(" ", child.w.value)
				}
			}

			if monospaceLength(line) < width {
				line += strings.Repeat(" ", width-monospaceLength(line))
			}
			fmt.Printf("adding line (col) %v\n", strings.ReplaceAll(line, " ", "~"))

			lines = append(lines, line)
		}
	}

	return lines
}

func (c node) String() string {
	return fmt.Sprintf("<width=%v>%v</>", c.w.value, c.content)
}

func monospaceLength(s string) int {
	return utf8.RuneCountInString(s)
}
