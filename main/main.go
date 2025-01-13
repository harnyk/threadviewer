package main

import (
	"fmt"
	"log"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sashabaranov/go-openai"
	"threadviewer/main/ui"
)

var (
	threadID string
	apiKey   string
)

var rootCmd = &cobra.Command{
	Use:   "threadviewer",
	Short: "threadviewer is a CLI tool to view threads from OpenAI",
	Run: func(cmd *cobra.Command, args []string) {
		client := openai.NewClient(apiKey)
		thread, err := client.RetrieveThread(threadID)
		if err != nil {
			log.Fatalf("Error retrieving thread: %v", err)
		}
		
		markdownContent := formatThreadAsMarkdown(thread)
		formattedOutput := ui.RenderMarkdown(markdownContent)
		fmt.Println(formattedOutput)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&threadID, "threadID", "", "Thread ID to retrieve")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "OpenAI API Key")
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("THREADVIEWER")
	viper.BindEnv("API_KEY")
	
	apiKey = viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}
}

func formatThreadAsMarkdown(thread openai.Thread) string {
	// Implement the formatting logic here
	return "# Thread Details\n\n" + thread.ID
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
