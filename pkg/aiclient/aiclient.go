package aiclient

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	*openai.Client
	Model       string
	Temperature float32
}

// Waits for the entire response to be returned
// Adds the users request, and the response to the conversation
func (c *Client) SendCompletionRequest(ctx context.Context, conv *Conversation, userPrompt string) (string, error) {
	// Add the latest message to the conversation
	err := conv.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userPrompt,
	})
	if err != nil {
		return "", err
	}

	// Send the request to the LLM 🤖
	completion, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       c.Model,
		Messages:    conv.Messages,
		Temperature: c.Temperature,
	})
	if err != nil {
		return "", err
	}
	responseChat := ""
	for _, token := range completion.Choices {
		responseChat = token.Message.Content
	}

	// Add the response to the conversation
	err = conv.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: responseChat,
	})
	if err != nil {
		return "", err
	}
	return responseChat, nil
}

// Stream the response as it comes in
// Adds the users request, and the response to the conversation
func (c *Client) SendStreamRequest(ctx context.Context, conv *Conversation, userPrompt string, responseChan chan string, errChan chan error) {
	defer close(responseChan)
	defer close(errChan)

	// Add the latest message to the conversation
	err := conv.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userPrompt,
	})
	if err != nil {
		errChan <- err
		return
	}

	// Stream the request to the LLM 🤖
	completionStream, err := c.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       c.Model,
		Temperature: c.Temperature,
		Messages:    conv.Messages,
		MaxTokens:   conv.maxTokens,
	})
	if err != nil {
		errChan <- err
		return
	}
	responseChat := ""
	for {
		streamData, err := completionStream.Recv()
		if err != nil {
			break
		}
		for _, token := range streamData.Choices {
			responseChan <- token.Delta.Content
			responseChat += token.Delta.Content
		}
	}
	// Add the response to the conversation, once the stream is closed
	err = conv.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: responseChat,
	})
	if err != nil {
		errChan <- err
		return
	}
	return
}
