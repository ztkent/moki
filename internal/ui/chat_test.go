package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	aiutil "github.com/ztkent/ai-util"
	"github.com/ztkent/moki/internal/chat"
	"github.com/ztkent/moki/internal/models"
)

// mockServer answers non-streaming requests with JSON and streaming requests
// with SSE, so a single server can drive both the intro and chat turns.
func mockServer(t *testing.T, reply string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Stream bool `json:"stream"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)

		if !body.Stream {
			fmt.Fprintf(w, `{"choices":[{"message":{"role":"assistant","content":%q},"finish_reason":"stop"}]}`, reply)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		fmt.Fprintf(w, "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\",\"content\":%q}}]}\n\n", reply)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":10,\"completion_tokens\":5,\"total_tokens\":15}}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
}

func newTestModel(t *testing.T, reply string) (ChatModel, *httptest.Server) {
	t.Helper()
	srv := mockServer(t, reply)
	client := aiutil.New("k", aiutil.WithBaseURL(srv.URL))
	session := chat.New(client, "test/model", "system", 0.5, 100)
	catalog := &models.Catalog{Models: []models.Model{{ID: "a/one", Name: "One"}}}
	return NewChatModel(context.Background(), session, catalog, "introduce yourself", "", nil), srv
}

func TestChatIntro(t *testing.T) {
	m, srv := newTestModel(t, "hello, I am Moki")
	defer srv.Close()

	updated, _ := m.Update(introMsg{reply: &chat.Reply{Text: "hello, I am Moki"}})
	m = updated.(ChatModel)
	if len(m.transcript) != 1 || !strings.Contains(m.transcript[0], "hello, I am Moki") {
		t.Errorf("transcript = %v", m.transcript)
	}
}

func TestChatUserMessageStreams(t *testing.T) {
	m, srv := newTestModel(t, "the answer")
	defer srv.Close()

	// Type a message and press enter.
	m.input.SetValue("what is 2+2?")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(ChatModel)

	if !m.streaming {
		t.Fatal("expected streaming to start")
	}
	if len(m.transcript) == 0 || !strings.Contains(m.transcript[0], "what is 2+2?") {
		t.Errorf("user message not recorded: %v", m.transcript)
	}

	// Run the stream command and feed the result back.
	msg := cmd()
	updated, _ = m.Update(msg)
	m = updated.(ChatModel)
	if m.streaming {
		t.Error("streaming should have finished")
	}
	// transcript: user, assistant reply (no per-turn footer).
	if len(m.transcript) != 2 {
		t.Fatalf("transcript length = %d, want 2: %v", len(m.transcript), m.transcript)
	}
	if !strings.Contains(m.transcript[1], "the answer") {
		t.Errorf("assistant reply not recorded: %v", m.transcript)
	}
	// The status line carries the cumulative token total.
	if !strings.Contains(m.statusLine(), "tokens") {
		t.Errorf("status line missing token total: %q", m.statusLine())
	}
}

func TestChatCommands(t *testing.T) {
	m, srv := newTestModel(t, "ok")
	defer srv.Close()

	// /help
	updated, _ := m.handleInput("/help")
	m = updated.(ChatModel)
	if !strings.Contains(strings.Join(m.transcript, "\n"), "/model") {
		t.Error("help text not shown")
	}

	// /model <id> switches the model.
	updated, _ = m.handleInput("/model b/two")
	m = updated.(ChatModel)
	if m.session.Model() != "b/two" {
		t.Errorf("model = %q, want b/two", m.session.Model())
	}

	// /clear empties the transcript.
	updated, _ = m.handleInput("/clear")
	m = updated.(ChatModel)
	if len(m.transcript) != 0 {
		t.Errorf("transcript not cleared: %v", m.transcript)
	}

	// /exit quits.
	updated, cmd := m.handleInput("/exit")
	m = updated.(ChatModel)
	if !m.quitting {
		t.Error("expected quitting")
	}
	if cmd == nil {
		t.Error("expected a quit command")
	}
}

func TestChatModelPicker(t *testing.T) {
	m, srv := newTestModel(t, "ok")
	defer srv.Close()

	// Size the model so the embedded picker gets dimensions.
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(ChatModel)

	// /model with no argument opens the picker.
	updated, _ = m.handleInput("/model")
	m = updated.(ChatModel)
	if m.picker == nil {
		t.Fatal("expected the picker to open")
	}
	// The embedded picker must be sized, or it renders nothing.
	if m.picker.list.Width() == 0 || m.picker.list.Height() == 0 {
		t.Errorf("picker not sized: %dx%d", m.picker.list.Width(), m.picker.list.Height())
	}
	if len(m.picker.list.Items()) == 0 {
		t.Error("picker has no items")
	}

	// Selecting an item applies it and closes the picker without quitting.
	chosen := models.Model{ID: "a/one", Name: "One"}
	pm := *m.picker
	pm.chosen = &chosen
	m.picker = &pm
	updated, cmd := m.updatePicker(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(ChatModel)
	if m.picker != nil {
		t.Error("picker should close after selection")
	}
	if m.session.Model() != "a/one" {
		t.Errorf("model = %q, want a/one", m.session.Model())
	}
	if m.quitting {
		t.Error("selecting a model must not quit the conversation")
	}
	if cmd != nil {
		t.Error("embedded picker should not return a quit command")
	}
}
