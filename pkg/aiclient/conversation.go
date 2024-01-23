package aiclient

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type Conversation struct {
	Messages []openai.ChatCompletionMessage
	lock     *sync.Mutex
	id       uuid.UUID
}

const (
	MAX_CONVERSATION_LENGTH = 100
)

// Start a new conversation with the system prompt
// A system prompt defines the initial context of the conversation
// This includes the persona of the bot and any information that you want to provide to the model.
func NewConversation(systemPrompt string) *Conversation {
	return &Conversation{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
		},
		lock: &sync.Mutex{},
		id:   uuid.New(),
	}
}

func (c *Conversation) Append(m openai.ChatCompletionMessage) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if len(c.Messages) >= MAX_CONVERSATION_LENGTH {
		return fmt.Errorf("conversation is at max length of %d", MAX_CONVERSATION_LENGTH)
	}
	c.Messages = append(c.Messages, m)
	return nil
}
