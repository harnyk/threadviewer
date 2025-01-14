package thread_service

import (
	"context"

	"github.com/harnyk/threadviewer/internal/hlp"
	"github.com/sashabaranov/go-openai"
)

type ThreadInfo struct {
	Messages            []openai.Message
	RunStepListsByRunID map[string]openai.RunStepList
	AssistantsByID      map[string]openai.Assistant
}

type ThreadService struct {
	client *openai.Client
}

func NewThreadService(client *openai.Client) *ThreadService {
	return &ThreadService{
		client: client,
	}
}

func (t *ThreadService) GetThreadInfo(ctx context.Context, threadID string) (*ThreadInfo, error) {
	info := &ThreadInfo{}

	threadMessages, err := t.client.ListMessage(
		ctx, threadID,
		hlp.PtrTo(100), hlp.PtrTo("asc"), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	info.Messages = threadMessages.Messages

	if err := t.getRunStepLists(ctx, info, threadID); err != nil {
		return nil, err
	}

	if err := t.getAssistants(ctx, info); err != nil {
		return nil, err
	}

	return info, nil
}

func (t *ThreadService) getAssistants(ctx context.Context, threadInfo *ThreadInfo) error {
	assistantsByID := make(map[string]openai.Assistant)
	for _, message := range threadInfo.Messages {
		if message.AssistantID != nil {
			assistantID := *message.AssistantID
			if _, ok := assistantsByID[assistantID]; ok {
				continue
			}

			assistant, err := t.client.RetrieveAssistant(ctx, assistantID)
			if err != nil {
				return err
			}

			assistantsByID[assistantID] = assistant
		}
	}

	threadInfo.AssistantsByID = assistantsByID
	return nil
}

func (t *ThreadService) getRunStepLists(ctx context.Context, threadInfo *ThreadInfo, threadID string) error {
	runStepListsByID := make(map[string]openai.RunStepList)
	for _, message := range threadInfo.Messages {
		if message.RunID != nil {
			runID := *message.RunID
			if _, ok := runStepListsByID[runID]; ok {
				continue
			}

			runSteps, err := t.client.ListRunSteps(ctx, threadID, runID, openai.Pagination{
				Limit: hlp.PtrTo(100),
				Order: hlp.PtrTo("asc"),
			})
			if err != nil {
				return err
			}

			runStepListsByID[runID] = runSteps
		}
	}
	threadInfo.RunStepListsByRunID = runStepListsByID

	return nil
}
