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