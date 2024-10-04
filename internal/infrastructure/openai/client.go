package openai

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/AsaHero/just-ask-bot/pkg/config"
	"github.com/sashabaranov/go-openai"
)

type apiClient struct {
	client *openai.Client
}

func New(cfg *config.Config) (OpenaiAPI, error) {
	client := openai.NewClient(cfg.OpenAI.SecretKey)

	return &apiClient{
		client: client,
	}, nil
}

func (c *apiClient) ChatCompletion(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (string, error) {
	var requestMessages []openai.ChatCompletionMessage
	requestMessages = append(requestMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: system,
	})

	requestMessages = append(requestMessages, messages...)

	request := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini20240718,
		Messages: requestMessages,
	}

	// Call the Chat Completion API
	response, err := c.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}

func (c *apiClient) ChatCompletionStreaming(ctx context.Context, system string, messages []openai.ChatCompletionMessage) (<-chan string, <-chan error) {
	output := make(chan string, 1)
	errchan := make(chan error, 1)
	chatCompletionMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: system,
		},
	}

	chatCompletionMessages = append(chatCompletionMessages, messages...)

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Stream:   true,
		Messages: chatCompletionMessages,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		close(output)
		errchan <- fmt.Errorf("(chatgptAPI) ChatCompletionStream error: %v", err)
		close(errchan)
		return output, errchan
	}

	go func() {
		defer close(output)
		defer close(errchan)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := stream.Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}

					errchan <- fmt.Errorf("(chatgptAPI), Error while parsing response chunks: %v", err)
					return
				}

				text := response.Choices[0].Delta.Content
				output <- text
			}
		}
	}()

	return output, nil
}
