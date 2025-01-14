package client

import (
	"context"

	"github.com/harnyk/threadviewer/internal/hlp"
	"github.com/sashabaranov/go-openai"
)

type ThreadInfo struct {
	Messages            []openai.Message
	RunStepListsByRunID map[string]openai.RunStepList
}

func GetThreadInfo(ctx context.Context, client *openai.Client, threadID string) (*ThreadInfo, error) {
	info := &ThreadInfo{}

	threadMessages, err := client.ListMessage(
		ctx, threadID,
		hlp.PtrTo(100), hlp.PtrTo("asc"), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	info.Messages = threadMessages.Messages

	rslsByID := make(map[string]openai.RunStepList)
	for _, message := range info.Messages {
		if message.RunID != nil {
			runSteps, err := client.ListRunSteps(ctx, threadID, *message.RunID, openai.Pagination{
				Limit: hlp.PtrTo(100),
				Order: hlp.PtrTo("asc"),
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
