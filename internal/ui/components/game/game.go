package game

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

var containerStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Margin(1, 1, 2, 1)

type Model struct {
	content  string
	viewport viewport.Model
}

func New() Model {
	return Model{
		content:  "test",
		viewport: viewport.New(),
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
	m.viewport.SetContent(m.content)
	return m.viewport.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case schedule.GameSelectedMsg:
		m.content = buildHeader(msg)
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func buildHeader(data schedule.GameSelectedMsg) string {
	header := data.Bug.AwayAbbr + " @ " + data.Bug.HomeAbbr

	return headerStyle.Render(header)
}
