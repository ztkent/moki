package conversation

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type MokiModel struct {
	textinput.Model
	quit              bool
	selectingResource bool
}

func (m MokiModel) Init() tea.Cmd {
	return nil
}

func (m MokiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "\x1b":
			m.quit = true
			return m, tea.Quit
		case "enter", "\r":
			return m, tea.Quit
		case "@":
			// If we are going to enter a resource, clear the view and reinvoke the text input
			if !m.selectingResource {
				m.selectingResource = true
				m.Blur()
				return m, tea.Tick(time.Millisecond, func(time.Time) tea.Msg {
					return tea.KeyMsg{
						Type:  tea.KeyRunes,
						Runes: []rune{'@'},
					}
				})
			}
			// With the input hidden, we can manage the resource selection
			modifiedInput, err := ManageResourceSelection(m.Value())
			m.selectingResource = false
			m.Focus()
			if err != nil {
				return m, tea.Quit
			} else if modifiedInput == m.Value() {
				// If the user cancels the resource selection, just return
				return m, nil
			}
			// Update the model with the modified input, including the resource
			m.SetValue(modifiedInput)
			return m, nil
		default:
			// Let the text input handle all other key presses
			updatedModel, cmd := m.Model.Update(msg)
			m.Model = updatedModel
			return m, cmd
		}
	}
	return m, nil
}

func (m MokiModel) View() string {
	if m.selectingResource {
		return ""
	}
	m.Focus()
	return m.Model.View()
}
