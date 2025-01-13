package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harnyk/threadviwer/internal/ui"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		ctx := context.Background()

		threadInfo, err := getThreadInfo(ctx, client, threadID)
		if err != nil {
			log.Fatalf("Error retrieving thread: %v", err)
		}

		markdownContent := formatThreadAsMarkdown(threadInfo)
		formattedOutput := ui.RenderMarkdown(markdownContent)
		fmt.Println(formattedOutput)
	},
}

type ThreadInfo struct {
	Messages            []openai.Message
	RunStepListsByRunID map[string]openai.RunStepList
}

func getThreadInfo(ctx context.Context, client *openai.Client, threadID string) (*ThreadInfo, error) {
	info := &ThreadInfo{}

	threadMessages, err := client.ListMessage(
		ctx, threadID,
		ptrTo(100), ptrTo("asc"), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	info.Messages = threadMessages.Messages

	rslsByID := make(map[string]openai.RunStepList)
	for _, message := range info.Messages {
		if message.RunID != nil {
			runSteps, err := client.ListRunSteps(ctx, threadID, *message.RunID, openai.Pagination{
				Limit: ptrTo(100),
				Order: ptrTo("asc"),
			})

			if err != nil {
				return nil, err
			}

			rslsByID[*message.RunID] = runSteps
		}
	}
	info.RunStepListsByRunID = rslsByID

	return info, nil
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&threadID, "threadID", "t", "", "Thread ID to retrieve")
	rootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "OpenAI API Key")
}

func initConfig() {
	viper.AutomaticEnv()
	// viper.SetEnvPrefix("THREADVIEWER")
	viper.BindEnv("API_KEY")

	apiKey = viper.GetString("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}
}

func formatThreadAsMarkdown(threadInfo *ThreadInfo) string {
	content := strings.Builder{}
	for _, message := range threadInfo.Messages {
		if message.RunID != nil {
			runStepList := threadInfo.RunStepListsByRunID[*message.RunID]
			content.WriteString(formatRunStepList(runStepList))
		}
		content.WriteString(fmt.Sprintf("%s **%s**: %s\n", getRoleEmoji(message.Role), message.Role, formatMessageContentList(message.Content)))
		content.WriteString("\n---\n")

	}
	return content.String()
}

func formatMessageContentList(content []openai.MessageContent) string {
	formattedContent := strings.Builder{}
	for _, messageContent := range content {
		formattedContent.WriteString(formatMessageContent(messageContent))
	}
	return formattedContent.String()
}

func formatRunStepList(runStepList openai.RunStepList) string {
	formattedRunSteps := strings.Builder{}
	for _, runStep := range runStepList.RunSteps {
		formattedRunSteps.WriteString(formatRunStep(runStep))
	}
	return formattedRunSteps.String()
}

func formatRunStep(runStep openai.RunStep) string {
	return formatRunToolCalls(runStep.StepDetails.ToolCalls)
}

func formatRunToolCalls(toolCalls []openai.ToolCall) string {
	formattedToolCalls := strings.Builder{}
	hammerEmoji := "🔨"
	for _, toolCall := range toolCalls {
		formattedToolCalls.WriteString(fmt.Sprintf("%s **%s(**%s**)**\n\n", hammerEmoji, toolCall.Function.Name, renderTextMaybeJSON(toolCall.Function.Arguments)))
		formattedToolCalls.WriteString(fmt.Sprintf("Returned: %s", renderTextMaybeJSON(toolCall.Function.Output)))
	}
	return formattedToolCalls.String()
}

func tryParseAsJSON(value string) (map[string]any, error) {
	var parsedJSON map[string]any
	if err := json.Unmarshal([]byte(value), &parsedJSON); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return parsedJSON, nil
}

func renderTextMaybeJSON(value string) string {
	parsedJSON, err := tryParseAsJSON(value)
	if err != nil {
		return value
	}
	encodedJSON, err := json.MarshalIndent(parsedJSON, "", "  ")
	if err != nil {
		return value
	}
	return "\n```json\n" + string(encodedJSON) + "\n```\n"
}

func formatMessageContent(content openai.MessageContent) string {
	switch content.Type {
	case "text":
		return renderTextMaybeJSON(content.Text.Value)
	case "image":
		return content.ImageFile.FileID
	default:
		return "Unsupported message type"
	}
}

func getRoleEmoji(role string) string {
	switch role {
	case "user":
		return "👤"
	case "assistant":
		return "🤖"
	default:
		return "🔲"
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ptrTo[T any](value T) *T {
	return &value
}
