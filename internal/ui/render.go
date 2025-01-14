package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/harnyk/threadviewer/internal/thread_service"
	"github.com/sashabaranov/go-openai"
)

type emoji string

const (
	hammerEmoji      emoji = "🔨"
	robotEmoji       emoji = "🤖"
	humanEmoji       emoji = "👤"
	whiteSquareEmoji emoji = "🔲"
)

func RenderThread(threadInfo *thread_service.ThreadInfo) string {
	content := strings.Builder{}
	for _, message := range threadInfo.Messages {
		if message.RunID != nil {
			runStepList := threadInfo.RunStepListsByRunID[*message.RunID]
			content.WriteString(renderRunStepList(runStepList))
		}

		content.WriteString(fmt.Sprintf("%s **%s**: %s\n",
			getRoleEmoji(message.Role),
			getRoleLabel(message, threadInfo),
			renderMessageContentList(message.Content),
		))
		content.WriteString("\n---\n")

	}
	return content.String()
}

func getAssistantName(message openai.Message, threadInfo *thread_service.ThreadInfo) string {
	if message.AssistantID == nil {
		return ""
	}

	assistantID := *message.AssistantID
	assistant, ok := threadInfo.AssistantsByID[assistantID]
	if !ok {
		return ""
	}

	namePtr := assistant.Name

	if namePtr == nil {
		return ""
	}

	return *namePtr
}

func renderMessageContentList(content []openai.MessageContent) string {
	formattedContent := strings.Builder{}
	for _, messageContent := range content {
		formattedContent.WriteString(renderMessageContent(messageContent))
	}
	return formattedContent.String()
}

func renderRunStepList(runStepList openai.RunStepList) string {
	formattedRunSteps := strings.Builder{}
	for _, runStep := range runStepList.RunSteps {
		formattedRunSteps.WriteString(renderRunStep(runStep))
	}
	return formattedRunSteps.String()
}

func renderRunStep(runStep openai.RunStep) string {
	return renderRunToolCalls(runStep.StepDetails.ToolCalls)
}

func renderRunToolCalls(toolCalls []openai.ToolCall) string {
	formattedToolCalls := strings.Builder{}
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

func renderMessageContent(content openai.MessageContent) string {
	switch content.Type {
	case "text":
		return renderTextMaybeJSON(content.Text.Value)
	case "image":
		return content.ImageFile.FileID
	default:
		return "Unsupported message type"
	}
}

func getRoleLabel(message openai.Message, threadInfo *thread_service.ThreadInfo) string {
	switch message.Role {
	case "user":
		return "User"
	case "assistant":
		if message.AssistantID != nil {
			asstName := getAssistantName(message, threadInfo)
			return fmt.Sprintf("Assistant (%s)", asstName)
		}
		return "Assistant"
	default:
		return message.Role
	}
}

func getRoleEmoji(role string) emoji {
	switch role {
	case "user":
		return humanEmoji
	case "assistant":
		return robotEmoji
	default:
		return whiteSquareEmoji
	}
}
