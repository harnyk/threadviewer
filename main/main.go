package main

import (
	"fmt"
	"os"
	"github.com/sashabaranov/go-openai"
	"ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: threadviewer <thread_id>")
		return
	}

	threadID := os.Args[1]
	client := openai.NewClient("your_api_key")

	thread, err := client.GetThread(threadID)
	if err != nil {
		fmt.Printf("Error retrieving thread: %v\n", err)
		return
	}

	markdownContent := formatThreadAsMarkdown(thread)
	formattedOutput := ui.RenderMarkdown(markdownContent)
	fmt.Println(formattedOutput)
}

func formatThreadAsMarkdown(thread openai.Thread) string {
	// Implement the formatting logic here
	return "# Thread Details\n\n" + thread.ID
}