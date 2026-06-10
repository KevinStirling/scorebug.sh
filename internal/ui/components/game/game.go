package game

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

type Model struct {
	game     *schedule.GameSelectedMsg
	viewport viewport.Model
}

func New() Model {
	return Model{
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
	m.viewport.SetContent(m.buildContent(m.viewport.Width()))
	return m.viewport.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case schedule.GameSelectedMsg:
		m.game = &msg
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) buildContent(width int) string {
	if m.game == nil {
		return ""
	}
	header := m.game.Bug.AwayAbbr + " @ " + m.game.Bug.HomeAbbr

	return headerStyle.Width(width).Render(header)
}
