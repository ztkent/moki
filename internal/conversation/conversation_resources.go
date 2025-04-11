package conversation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ResourceSelectionModel struct {
	resourceTypes []string
	cursor        int
	selected      bool
	quit          bool
	textinput.Model
}

func (m ResourceSelectionModel) Init() tea.Cmd {
	return nil
}
func (m ResourceSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "esc", "\x1b":
			m.quit = true
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
	if m.selected || m.quit {
		return ""
	}
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

type ResourceInputModel struct {
	textinput.Model
	inputType string
	finished  bool
}

func (m ResourceInputModel) Init() tea.Cmd {
	return nil
}

func (m ResourceInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "\x1b":
			return m, tea.Quit
		case "enter", "\r":
			m.finished = true
			return m, tea.Quit
		default:
			// Let the text input handle all other key presses
			updatedModel, cmd := m.Model.Update(msg)
			m.Model = updatedModel
			return m, cmd
		}
	}
	return m, nil
}

func (m ResourceInputModel) View() string {
	if m.finished {
		return ""
	}
	return m.Model.View()
}

func ManageResourceSelection(userInput string) (string, error) {
	// Select the type of resource to input
	resourceType, err := getResourceType(userInput)
	if err != nil {
		return userInput, err
	} else if resourceType == "" {
		return userInput, nil
	}

	// Get the resource path
	resourcePath, err := getResourcePath(resourceType)
	if err != nil {
		return userInput, err
	} else if resourcePath == "" {
		return userInput, nil
	}

	// Add the resource to the user's input
	return userInput + " -" + resourceType + ":" + resourcePath, nil
}

func getResourcePath(resourceType string) (string, error) {
	m := ResourceInputModel{inputType: resourceType, Model: textinput.New()}
	m.Focus()
	defer m.Blur()
	m.Prompt = resourceType + ": "
	p := tea.NewProgram(m)
	defer p.RestoreTerminal()
	if resModel, err := p.Run(); err != nil {
		return "", err
	} else {
		if !resModel.(ResourceInputModel).finished {
			return "", nil
		}
		m = resModel.(ResourceInputModel)
	}
	return m.Value(), nil
}

func getResourceType(userInput string) (string, error) {
	m := ResourceSelectionModel{resourceTypes: []string{"url", "file"}, Model: textinput.New()}
	m.Focus()
	defer m.Blur()
	p := tea.NewProgram(m)
	defer p.RestoreTerminal()
	if resModel, err := p.Run(); err != nil {
		return userInput, err
	} else {
		if !resModel.(ResourceSelectionModel).selected {
			return "", nil
		}
		m = resModel.(ResourceSelectionModel)
	}
	return m.resourceTypes[m.cursor], nil
}
