package aiclient

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	Messages    []openai.ChatCompletionMessage
	maxMessages int
	maxTokens   int
	tokenCount  int
	lock        *sync.Mutex
	id          uuid.UUID
}

const (
	MAX_CONVERSATION_LENGTH = 100
	MAX_TOKENS              = 1600
)

// Start a new conversation with the system prompt
// A system prompt defines the initial context of the conversation
// This includes the persona of the bot and any information that you want to provide to the model.
func NewConversation(systemPrompt string, maxMessages int, maxTokens int) *Conversation {
	if maxMessages == 0 {
		maxMessages = MAX_CONVERSATION_LENGTH
	}
	if maxTokens == 0 {
		maxTokens = MAX_TOKENS
	}
	conv := &Conversation{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
		},
		maxMessages: maxMessages,
		maxTokens:   maxTokens,
		lock:        &sync.Mutex{},
		id:          uuid.New(),
	}
	return conv
}

func (c *Conversation) Append(m openai.ChatCompletionMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if len(c.Messages) >= c.maxMessages {
		return fmt.Errorf("conversation is at max length of %d", c.maxMessages)
	}
	c.Messages = append(c.Messages, m)
	return nil
}

func (c *Conversation) SeedConversation(requestResponseMap map[string]string) {
	// Seed the conversation with some example prompts and responses
	for user, response := range requestResponseMap {
		c.Append(openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: user,
		})
		c.Append(openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: response,
		})
	}
}
