package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	aiutil "github.com/ztkent/ai-util"
)

// streamServer replies with a fixed streamed message.
func streamServer(t *testing.T, reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\",\"content\":%q}}]}\n\n", reply)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
}

func TestSendStream(t *testing.T) {
	srv := streamServer(t, "hello there")
	defer srv.Close()

	client := aiutil.New("k", aiutil.WithBaseURL(srv.URL))
	s := New(client, "test/model", "be helpful", 0.5, 100, 0)

	var streamed strings.Builder
	reply, err := s.SendStream(context.Background(), "hi", func(text string) {
		streamed.WriteString(text)
	})
	if err != nil {
		t.Fatalf("SendStream: %v", err)
	}
	if reply.Text != "hello there" {
		t.Errorf("reply = %q, want %q", reply.Text, "hello there")
	}
	if streamed.String() != "hello there" {
		t.Errorf("streamed = %q, want %q", streamed.String(), "hello there")
	}
}

func TestSetModelAndReset(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`)
	}))
	defer srv.Close()

	client := aiutil.New("k", aiutil.WithBaseURL(srv.URL))
	s := New(client, "model-a", "", 0.5, 100, 0)

	if s.Model() != "model-a" {
		t.Errorf("model = %q, want model-a", s.Model())
	}
	s.SetModel("model-b")
	if s.Model() != "model-b" {
		t.Errorf("model = %q, want model-b", s.Model())
	}

	if _, err := s.Send(context.Background(), "hi"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if len(s.agent.History()) == 0 {
		t.Fatal("expected history to be populated")
	}
	s.Reset()
	if len(s.agent.History()) != 0 {
		t.Errorf("history length = %d after reset, want 0", len(s.agent.History()))
	}
}

func TestRequestIncludesSystemPrompt(t *testing.T) {
	s := New(aiutil.New("k"), "m", "system text", 0.3, 50, 0)
	if s.system != "system text" {
		t.Errorf("system = %q, want %q", s.system, "system text")
	}
}

// TestSetModelKeepsHistory verifies switching models preserves the conversation.
func TestSetModelKeepsHistory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`)
	}))
	defer srv.Close()

	s := New(aiutil.New("k", aiutil.WithBaseURL(srv.URL)), "model-a", "", 0.5, 100, 0)
	if _, err := s.Send(context.Background(), "hi"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	before := len(s.agent.History())

	s.SetModel("model-b")
	if got := len(s.agent.History()); got != before {
		t.Errorf("history length = %d after SetModel, want %d", got, before)
	}
}

// TestParamsReachWire verifies the session's temperature and max tokens are
// injected into every request, since the agent builds its own.
func TestParamsReachWire(t *testing.T) {
	var got struct {
		Temperature *float64 `json:"temperature"`
		MaxTokens   int      `json:"max_tokens"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&got)
		fmt.Fprint(w, `{"choices":[{"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`)
	}))
	defer srv.Close()

	s := New(aiutil.New("k", aiutil.WithBaseURL(srv.URL)), "m", "", 0.3, 50, 0)
	if _, err := s.Send(context.Background(), "hi"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got.Temperature == nil || *got.Temperature != 0.3 {
		t.Errorf("temperature = %v, want 0.3", got.Temperature)
	}
	if got.MaxTokens != 50 {
		t.Errorf("max_tokens = %d, want 50", got.MaxTokens)
	}
}
