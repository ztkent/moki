package chat

import (
	"context"
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
	s := New(client, "test/model", "be helpful", 0.5, 100)

	var streamed strings.Builder
	got, err := s.SendStream(context.Background(), "hi", func(text string) {
		streamed.WriteString(text)
	})
	if err != nil {
		t.Fatalf("SendStream: %v", err)
	}
	if got != "hello there" {
		t.Errorf("reply = %q, want %q", got, "hello there")
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
	s := New(client, "model-a", "", 0.5, 100)

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
	if len(s.history) == 0 {
		t.Fatal("expected history to be populated")
	}
	s.Reset()
	if len(s.history) != 0 {
		t.Errorf("history length = %d after reset, want 0", len(s.history))
	}
}

func TestRequestIncludesSystemPrompt(t *testing.T) {
	s := New(aiutil.New("k"), "m", "system text", 0.3, 50)
	req := s.request()
	if len(req.Messages) != 1 || req.Messages[0].Role != aiutil.RoleSystem {
		t.Fatalf("expected a single system message, got %+v", req.Messages)
	}
	if req.Model != "m" || req.MaxTokens != 50 {
		t.Errorf("unexpected request: %+v", req)
	}
	if req.Temperature == nil || *req.Temperature != 0.3 {
		t.Errorf("temperature = %v, want 0.3", req.Temperature)
	}
}
