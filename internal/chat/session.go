// Package chat wraps an ai-util client with conversation history.
package chat

import (
	"context"

	aiutil "github.com/ztkent/ai-util"
)

// Session is a stateful chat with a single model.
type Session struct {
	client      *aiutil.Client
	model       string
	system      string
	temperature *float64
	maxTokens   int
	history     []aiutil.Message
}

// New creates a Session.
func New(client *aiutil.Client, model, system string, temperature float64, maxTokens int) *Session {
	t := temperature
	return &Session{
		client:      client,
		model:       model,
		system:      system,
		temperature: &t,
		maxTokens:   maxTokens,
	}
}

// Model returns the model currently in use.
func (s *Session) Model() string { return s.model }

// SetModel switches the model for subsequent turns, keeping the history.
func (s *Session) SetModel(model string) { s.model = model }

// Reset clears the conversation history.
func (s *Session) Reset() { s.history = nil }

// Send appends input, gets a complete reply, and returns its text.
func (s *Session) Send(ctx context.Context, input string) (string, error) {
	s.history = append(s.history, aiutil.User(input))
	resp, err := s.client.Chat(ctx, s.request())
	if err != nil {
		return "", err
	}
	s.history = append(s.history, resp.Message)
	return resp.Message.Content, nil
}

// SendStream is Send with each text chunk delivered to onText.
func (s *Session) SendStream(ctx context.Context, input string, onText func(string)) (string, error) {
	s.history = append(s.history, aiutil.User(input))
	resp, err := s.client.ChatStream(ctx, s.request(), func(e aiutil.Event) error {
		if e.Type == aiutil.EventText && onText != nil {
			onText(e.Text)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	s.history = append(s.history, resp.Message)
	return resp.Message.Content, nil
}

// request builds the wire request from the system prompt and history.
func (s *Session) request() *aiutil.Request {
	msgs := make([]aiutil.Message, 0, len(s.history)+1)
	if s.system != "" {
		msgs = append(msgs, aiutil.System(s.system))
	}
	msgs = append(msgs, s.history...)

	return &aiutil.Request{
		Model:       s.model,
		Messages:    msgs,
		Temperature: s.temperature,
		MaxTokens:   s.maxTokens,
	}
}
