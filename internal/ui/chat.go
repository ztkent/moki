package ui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ztkent/moki/internal/chat"
	"github.com/ztkent/moki/internal/models"
)

const chatHelp = `Commands:
  /model [id]   Switch model (opens the picker when no id is given)
  /clear        Clear the conversation history
  /help         Show this help
  /exit         Quit (also /quit, /q, or ctrl+c)`

// streamMsg carries a chunk of streamed assistant text.
type streamMsg string

// streamDoneMsg signals the end of a streamed reply and carries its full text.
type streamDoneMsg struct {
	text string
	err  error
}

// introMsg carries Moki's opening greeting.
type introMsg struct {
	text string
	err  error
}

// ChatModel is the interactive conversation TUI.
type ChatModel struct {
	ctx         context.Context
	session     *chat.Session
	catalog     *models.Catalog
	introPrompt string

	viewport viewport.Model
	input    textinput.Model
	picker   *PickerModel

	transcript []string
	pending    string
	streaming  bool
	quitting   bool
	ready      bool

	program *tea.Program
}

// NewChatModel builds the conversation model.
func NewChatModel(ctx context.Context, session *chat.Session, catalog *models.Catalog, introPrompt string) ChatModel {
	input := textinput.New()
	input.Prompt = "You: "
	input.Placeholder = "ask anything, or /help"
	input.Focus()

	return ChatModel{
		ctx:         ctx,
		session:     session,
		catalog:     catalog,
		introPrompt: introPrompt,
		viewport:    viewport.New(0, 0),
		input:       input,
	}
}

// RunChat runs the conversation TUI until the user exits.
func RunChat(ctx context.Context, session *chat.Session, catalog *models.Catalog, introPrompt string) error {
	m := NewChatModel(ctx, session, catalog, introPrompt)
	p := tea.NewProgram(m)
	m.program = p
	_, err := p.Run()
	return err
}

func (m ChatModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.fetchIntro())
}

func (m ChatModel) fetchIntro() tea.Cmd {
	return func() tea.Msg {
		text, err := m.session.Send(m.ctx, m.introPrompt)
		return introMsg{text: text, err: err}
	}
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.viewport.Width = msg.Width
		m.viewport.Height = max(msg.Height-2, 1)
		m.input.Width = max(msg.Width-len("You: ")-1, 1)
		if m.picker != nil {
			updated, _ := m.picker.Update(msg)
			pm := updated.(PickerModel)
			m.picker = &pm
		}
		m.render()
		return m, nil

	case tea.KeyMsg:
		if m.picker != nil {
			return m.updatePicker(msg)
		}
		if m.streaming {
			if msg.String() == "ctrl+c" {
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			value := strings.TrimSpace(m.input.Value())
			if value == "" {
				return m, nil
			}
			m.input.SetValue("")
			return m.handleInput(value)
		}
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd

	case streamMsg:
		m.pending += string(msg)
		m.render()
		return m, nil

	case streamDoneMsg:
		m.streaming = false
		switch {
		case msg.err != nil:
			m.transcript = append(m.transcript, errorStyle.Render("error: "+msg.err.Error()))
		case msg.text != "":
			m.transcript = append(m.transcript, assistantBlock(msg.text))
		}
		m.pending = ""
		m.render()
		return m, nil

	case introMsg:
		if msg.err != nil {
			m.transcript = append(m.transcript, errorStyle.Render("error: "+msg.err.Error()))
		} else {
			m.transcript = append(m.transcript, assistantBlock(msg.text))
		}
		m.render()
		return m, nil
	}
	return m, nil
}

// updatePicker delegates input to the model picker and applies the result.
func (m ChatModel) updatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.picker.Update(msg)
	pm := updated.(PickerModel)
	m.picker = &pm

	switch {
	case pm.Chosen() != nil:
		m.session.SetModel(pm.Chosen().ID)
		m.transcript = append(m.transcript, systemBlock("model → "+pm.Chosen().ID))
		m.picker = nil
		m.render()
		return m, nil
	case pm.cancelled:
		m.picker = nil
		m.render()
		return m, nil
	}
	return m, cmd
}

// handleInput routes slash commands or starts a streamed reply.
func (m ChatModel) handleInput(value string) (tea.Model, tea.Cmd) {
	if strings.HasPrefix(value, "/") {
		return m.handleCommand(value)
	}
	m.transcript = append(m.transcript, userBlock(value))
	m.streaming = true
	m.pending = ""
	m.render()
	return m, m.startStream(value)
}

func (m ChatModel) handleCommand(value string) (tea.Model, tea.Cmd) {
	fields := strings.Fields(value)
	switch strings.ToLower(fields[0]) {
	case "/exit", "/quit", "/q":
		m.quitting = true
		return m, tea.Quit
	case "/clear":
		m.session.Reset()
		m.transcript = nil
		m.pending = ""
		m.render()
		return m, nil
	case "/model":
		if len(fields) > 1 {
			m.session.SetModel(fields[1])
			m.transcript = append(m.transcript, systemBlock("model → "+fields[1]))
			m.render()
			return m, nil
		}
		if m.catalog == nil {
			m.transcript = append(m.transcript, systemBlock("no model catalog available"))
			m.render()
			return m, nil
		}
		p := NewPicker(m.catalog)
		m.picker = &p
		return m, nil
	case "/help":
		m.transcript = append(m.transcript, systemBlock(chatHelp))
		m.render()
		return m, nil
	default:
		m.transcript = append(m.transcript, systemBlock("unknown command: "+fields[0]))
		m.render()
		return m, nil
	}
}

// startStream runs a streamed reply, forwarding chunks to the program for live
// display and returning the full text when the stream ends.
func (m ChatModel) startStream(prompt string) tea.Cmd {
	return func() tea.Msg {
		var full strings.Builder
		_, err := m.session.SendStream(m.ctx, prompt, func(text string) {
			full.WriteString(text)
			if m.program != nil {
				m.program.Send(streamMsg(text))
			}
		})
		return streamDoneMsg{text: full.String(), err: err}
	}
}

// render rebuilds the viewport from the transcript and in-flight reply.
func (m *ChatModel) render() {
	blocks := make([]string, 0, len(m.transcript)+1)
	blocks = append(blocks, m.transcript...)
	if m.pending != "" {
		blocks = append(blocks, assistantBlock(m.pending))
	}
	m.viewport.SetContent(strings.Join(blocks, "\n\n"))
	m.viewport.GotoBottom()
}

func (m ChatModel) View() string {
	if m.quitting {
		return ""
	}
	if m.picker != nil {
		return m.picker.View()
	}
	if !m.ready {
		return statusStyle.Render("Starting Moki…")
	}
	return m.viewport.View() + "\n" + m.input.View()
}

func userBlock(text string) string      { return userStyle.Render("You: ") + text }
func assistantBlock(text string) string { return mokiStyle.Render("Moki: ") + text }
func systemBlock(text string) string    { return helpStyle.Render(text) }
