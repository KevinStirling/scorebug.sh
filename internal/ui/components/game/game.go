package game

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var containerStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Foreground(lipgloss.Green).
	Margin(1, 1)

type Model struct {
	GameContent     string
	ContainerWidth  int
	ContainerHeight int
	viewport        viewport.Model
}

func New() Model {
	return Model{
		GameContent: "test",
		viewport:    viewport.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch()
}

func (m *Model) SetSize(width, height int) {
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
}

func (m Model) View() string {
	m.viewport.Style = containerStyle
	m.viewport.SetContent(m.GameContent)
	return m.viewport.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg.(type) {
	case Model:
		m.viewport.SetContent(m.GameContent)
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, tea.Batch(cmd)
}
