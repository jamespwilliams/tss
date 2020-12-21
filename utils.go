package tss

import (
	"strings"
	"unicode/utf8"
)

func monospaceLength(s string) int {
	return utf8.RuneCountInString(s)
}

func addBorder(lines []string, width int) []string {
	if len(lines) == 0 {
		return []string{}
	}

	// fmt.Println(width, len(lines[0]))

	res := []string{"┌" + strings.Repeat("─", width) + "┐"}
	// fmt.Println("boxing           ", res[0])

	for _, line := range lines {
		res = append(res, "│"+line+"│")
		// fmt.Println("boxing           ", res[len(res)-1])
	}

	res = append(res, "└"+strings.Repeat("─", width)+"┘")
	// fmt.Println("boxing           ", res[len(res)-1])
	return res
}
