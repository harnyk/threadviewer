package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/harnyk/threadviewer/internal/client"
	"github.com/harnyk/threadviewer/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	threadID   string
	apiKey     string
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "threadviewer <threadID>",
	Short: "threadviewer is a CLI tool to view threads from OpenAI",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("threadID is required")
		}

		threadID = args[0]
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		openAI := openai.NewClient(apiKey)
		ctx := context.Background()

		threadInfo, err := client.GetThreadInfo(ctx, openAI, threadID)
		if err != nil {
			log.Fatalf("Error retrieving thread: %v", err)
		}

		markdownContent := ui.RenderThread(threadInfo)
		formattedOutput := ui.RenderMarkdown(markdownContent)
		fmt.Println(formattedOutput)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "OpenAI API Key")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Config file path")
}

func initConfig() {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".threadviewer")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	viper.BindEnv("API_KEY")

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	apiKey = viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
