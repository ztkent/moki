package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ztkent/moki/internal/models"
)

// modelItem adapts a catalog entry to the bubbles list.
type modelItem struct {
	model     models.Model
	preferred bool
}

func (i modelItem) Title() string {
	if i.preferred {
		return "★ " + i.model.Name
	}
	return i.model.Name
}

func (i modelItem) FilterValue() string {
	return i.model.ID + " " + i.model.Name
}

func (i modelItem) Description() string {
	var parts []string
	parts = append(parts, i.model.ID)
	if i.preferred {
		parts = append(parts, "preferred")
	}
	if i.model.IsFree() {
		parts = append(parts, "free")
	}
	if i.model.ContextLength > 0 {
		parts = append(parts, fmt.Sprintf("%s ctx", humanTokens(i.model.ContextLength)))
	}
	return strings.Join(parts, " · ")
}

// humanTokens renders a token count like 200K or 1M.
func humanTokens(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.0fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.0fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// PickerModel is an interactive, filterable model chooser.
type PickerModel struct {
	list      list.Model
	chosen    *models.Model
	cancelled bool
}

// NewPicker builds a picker from a catalog. The model matching preferred is
// marked and pre-selected.
func NewPicker(catalog *models.Catalog, preferred string) PickerModel {
	items := make([]list.Item, 0, len(catalog.Models))
	selected := 0
	for i, m := range catalog.Models {
		isPreferred := m.ID == preferred
		if isPreferred {
			selected = i
		}
		items = append(items, modelItem{model: m, preferred: isPreferred})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(accent).BorderForeground(accent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(accent).BorderForeground(accent)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Choose a model"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(accent)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(accent)
	if preferred != "" {
		l.Select(selected)
	}

	return PickerModel{list: l}
}

// Chosen returns the selected model, or nil when the picker was cancelled.
func (m PickerModel) Chosen() *models.Model { return m.chosen }

func (m PickerModel) Init() tea.Cmd { return nil }

func (m PickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		// Don't let quit keys fire while the user is typing a filter.
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "ctrl+c", "esc":
				m.cancelled = true
				return m, tea.Quit
			case "enter":
				if item, ok := m.list.SelectedItem().(modelItem); ok {
					chosen := item.model
					m.chosen = &chosen
				}
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PickerModel) View() string {
	return m.list.View()
}

// RunPicker shows the picker and returns the chosen model, or nil if cancelled.
func RunPicker(catalog *models.Catalog, preferred string, out io.Writer) (*models.Model, error) {
	p := tea.NewProgram(NewPicker(catalog, preferred), tea.WithOutput(out))
	res, err := p.Run()
	if err != nil {
		return nil, err
	}
	return res.(PickerModel).Chosen(), nil
}
