# Task

You are experienced Go developer with AI background.

Create a Go command line application which would load a given assystant thread by ID and print it to the stdout, applying markdown formatting to the output.

The output should be formatted in such way that it would be suitable for rendering in a terminal and readable by a human.

You are creating the application from scratch.

The app is called threadviewer.

The main package is called `main`.

Notes:

For markdown rendering create and use the following function:

```go
package ui
import (
	"os"
	markdown "github.com/MichaelMure/go-term-markdown"
	"golang.org/x/term"
)
func RenderMarkdown(source string) string {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	return string(markdown.Render(source, width, 0))
}
```

For retrieving the thread use github.com/sashabaranov/go-openai.

