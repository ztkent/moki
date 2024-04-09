package conversation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MokiModel struct {
	textInput textinput.Model
	quit      bool
}

func NewMokiModel() MokiModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 2560
	ti.Width = 20
	return MokiModel{textInput: ti}
}

func (m MokiModel) Init() tea.Cmd {
	return nil
}

func (m MokiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quit = true
			return m, tea.Quit
		case "esc", "\x1b":
			m.quit = true
			return m, tea.Quit
		case "enter", "\r":
			return m, tea.Quit
		case "@":
			// Manage resource selection
			modifiedInput, err := ManageResourceSelection(m.textInput.Value())
			if err != nil {
				return m, tea.Quit
			}
			m.textInput.SetValue(modifiedInput)
		default:
			// Let the text input handle all other key presses
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m MokiModel) View() string {
	return "You: " + m.textInput.Value()
}

type ResourceSelectionModel struct {
	resourceTypes []string
	cursor        int
	selected      bool
}

func (m ResourceSelectionModel) Init() tea.Cmd {
	return nil
}
func (m ResourceSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q", "\x03":
			return m, tea.Quit
		case "esc", "\x1b":
			return m, tea.Quit
		case "enter", "\r":
			m.selected = true
			return m, tea.Quit
		case "down", "\x1b[B":
			if m.cursor < len(m.resourceTypes)-1 {
				m.cursor++
			}
		case "up", "\x1b[A":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m ResourceSelectionModel) View() string {
	view := ""
	for i, resourceType := range m.resourceTypes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		view += fmt.Sprintf("%s %s\n", cursor, resourceType)
	}
	return view
}

func ManageResourceSelection(userInput string) (string, error) {
	m := ResourceSelectionModel{resourceTypes: []string{"url", "file"}}
	p := tea.NewProgram(m)
	if m, err := p.Run(); err != nil {
		return userInput, err
	} else {
		if !m.(ResourceSelectionModel).selected {
			return userInput, nil
		}
	}
	resourceType := m.resourceTypes[m.cursor]

	// Prompt the user to enter the resource
	fmt.Print("Enter the " + resourceType + ": ")
	reader := bufio.NewReader(os.Stdin)
	resource, _ := reader.ReadString('\n')
	resource = strings.TrimSpace(resource)

	// Add the resource to the user's input
	return userInput + " -" + resourceType + ":" + resource, nil
}
