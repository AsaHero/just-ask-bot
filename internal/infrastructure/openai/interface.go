package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenaiAPI interface {
	ChatCompletion(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (string, error)
	ChatCompletionStreaming(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (<-chan string, <-chan error)
}
