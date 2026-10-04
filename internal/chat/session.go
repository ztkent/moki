// Package chat wraps an ai-util agent with conversation history.
package chat

import (
	"context"

	aiutil "github.com/ztkent/ai-util"
)

// Session is a stateful chat with a single model, backed by an ai-util agent.
type Session struct {
	client           aiutil.ChatClient
	model            string
	system           string
	maxContextTokens int
	agent            *aiutil.Agent
}

// New creates a Session. temperature and maxTokens apply to every model call;
// maxContextTokens trims the oldest history to stay within budget (0 disables).
func New(client aiutil.ChatClient, model, system string, temperature float64, maxTokens, maxContextTokens int) *Session {
	s := &Session{
		client:           paramClient{client, temperature, maxTokens},
		model:            model,
		system:           system,
		maxContextTokens: maxContextTokens,
	}
	s.agent = s.newAgent()
	return s
}

// newAgent builds an agent for the current model, seeded with history.
func (s *Session) newAgent(history ...aiutil.Message) *aiutil.Agent {
	return aiutil.NewAgent(s.client, s.model,
		aiutil.WithSystem(s.system),
		aiutil.WithMaxContextTokens(s.maxContextTokens),
		aiutil.WithHistory(history...),
	)
}

// Model returns the model currently in use.
func (s *Session) Model() string { return s.model }

// SetModel switches the model for subsequent turns, keeping the history.
func (s *Session) SetModel(model string) {
	s.model = model
	s.agent = s.newAgent(s.agent.History()...)
}

// Reset clears the conversation history.
func (s *Session) Reset() { s.agent.Reset() }

// Reply is the result of a turn, including which model answered and its usage.
type Reply struct {
	Text  string
	Model string
	Usage aiutil.Usage
}

// Send gets a complete reply to input.
func (s *Session) Send(ctx context.Context, input string) (*Reply, error) {
	resp, err := s.agent.Run(ctx, input)
	if err != nil {
		return nil, err
	}
	return s.replyFrom(resp), nil
}

// SendStream is Send with each text chunk delivered to onText.
func (s *Session) SendStream(ctx context.Context, input string, onText func(string)) (*Reply, error) {
	resp, err := s.agent.RunStream(ctx, input, func(e aiutil.Event) error {
		if e.Type == aiutil.EventText && onText != nil {
			onText(e.Text)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.replyFrom(resp), nil
}

// replyFrom converts an ai-util response into a Reply, falling back to the
// requested model when the provider omits it.
func (s *Session) replyFrom(resp *aiutil.Response) *Reply {
	model := resp.Model
	if model == "" {
		model = s.model
	}
	return &Reply{Text: resp.Message.Content, Model: model, Usage: resp.Usage}
}

// paramClient injects the session's sampling parameters into every request.
// ai-util's Agent builds its own requests, so this is how temperature and max
// tokens reach the wire.
type paramClient struct {
	client      aiutil.ChatClient
	temperature float64
	maxTokens   int
}

func (p paramClient) Chat(ctx context.Context, req *aiutil.Request) (*aiutil.Response, error) {
	p.apply(req)
	return p.client.Chat(ctx, req)
}

func (p paramClient) ChatStream(ctx context.Context, req *aiutil.Request, onEvent func(aiutil.Event) error) (*aiutil.Response, error) {
	p.apply(req)
	return p.client.ChatStream(ctx, req, onEvent)
}

func (p paramClient) apply(req *aiutil.Request) {
	t := p.temperature
	req.Temperature = &t
	req.MaxTokens = p.maxTokens
}
