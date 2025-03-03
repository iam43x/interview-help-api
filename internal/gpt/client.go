package gpt

import (
	"context"
	"fmt"

	"github.com/iam43x/interview-help-4u/internal/util"

	openai "github.com/sashabaranov/go-openai"
)

type ChatGptClient struct {
	OpenaiClient *openai.Client
}

func NewChatGptClient(apiKey string) *ChatGptClient {
	return &ChatGptClient{
		OpenaiClient: openai.NewClient(apiKey),
	}
}

func (c *ChatGptClient) TranscribeAudio(ctx context.Context, r *util.WriteSeeker) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		Reader:   r.Reader(),
		FilePath: r.Filename,
	}

	resp, err := c.OpenaiClient.CreateTranscription(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к OpenAI: %v", err)
	}

	return resp.Text, nil
}

func (c *ChatGptClient) AskGpt3Dot5Turbo16K(ctx context.Context, question string) (string, error) {
	resp, err := c.OpenaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo16K,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к OpenAI: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}