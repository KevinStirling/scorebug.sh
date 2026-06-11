package game

import (
	"strconv"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
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
	m.viewport.SetContent(m.renderContent(m.viewport.Width()))
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

func (m Model) renderContent(width int) string {
	if m.game == nil {
		return ""
	}
	header := m.game.Bug.Feed.GameData.Teams.Away.Name + " @ " + m.game.Bug.Feed.GameData.Teams.Home.Name
	linescore := m.buildLineScore()

	content := lipgloss.JoinVertical(lipgloss.Center, header, linescore)

	return headerStyle.Width(width).Render(content)
}

func (m Model) buildLineScore() string {
	inningRow := []string{""}
	awayRow := []string{m.game.Bug.AwayAbbr}
	homeRow := []string{m.game.Bug.HomeAbbr}
	linescore := m.game.Bug.Feed.LiveData.Linescore
	for _, v := range linescore.Innings {
		inningRow = append(inningRow, strconv.Itoa(v.Num))
		awayRow = append(awayRow, strconv.Itoa(v.Away.Runs))
		homeRow = append(homeRow, strconv.Itoa(v.Home.Runs))
	}
	inningRow = append(inningRow, []string{"R", "H", "E"}...)
	awayRow = append(awayRow, []string{strconv.Itoa(linescore.Teams.Away.Runs), strconv.Itoa(linescore.Teams.Away.Hits), strconv.Itoa(linescore.Teams.Away.Errors)}...)
	homeRow = append(homeRow, []string{strconv.Itoa(linescore.Teams.Home.Runs), strconv.Itoa(linescore.Teams.Home.Hits), strconv.Itoa(linescore.Teams.Home.Errors)}...)

	rows := [][]string{
		awayRow,
		homeRow,
	}

	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return cellStyle.Foreground(lipgloss.Magenta).Bold(true)
			}
			return cellStyle
		}).
		Headers(inningRow...).Rows(rows...)
	return t.Render()
}
