# tss - Terminal Style Sheets

Go library for writing TUIs in a HTML/CSS-esque syntax.

Heavily WIP.

## Example

```golang
package main

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jamespwilliams/tss"
)

func main() {
	s := strings.NewReader(`
		<div flow="column" width="100" border>
			<div flow="row" id="one" width="35" border>
				QQQQQQQQQQQQQQQQQ
				<div flow="column" border>
					<div width="10" border>
						ABC
					</div>
					<div width="10" border>
						XYZ
					</div>
				</div>
			</div>
			<div id="two" width="30" border>
			</div>
		</div>`)

	t, _ := tss.NewTemplate(s, nil)

	t.Render(os.Stdout)

	for {
		t.RenderFullScreen(os.Stdout)

		size := rand.Intn(25)
		elem := t.Elements["two"]
		elem.SetContent(randSeq(size))
		elem.SetWidth(size)

		time.Sleep(1 * time.Second)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
```

## TODOs

- Background colors
- Font colors
- Better layout options
    - Allow (equivalents to) flexbox's `justify-content` and `align-items`
- Text align
- Different border styles?
