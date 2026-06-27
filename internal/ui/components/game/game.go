package game

import (
	"strconv"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/internal/ui/components/schedule"
)

type Model struct {
	game     *schedule.GameSelectedMsg
	viewport viewport.Model
	plays    PlayFeed
}

func New() Model {
	return Model{
		viewport: viewport.New(),
		plays:    NewPlayfeed(),
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
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.plays, cmd = m.plays.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) renderContent(width int) string {
	if m.game == nil {
		return ""
	}
	header := m.game.Bug.Feed.GameData.Teams.Away.Name + " @ " + m.game.Bug.Feed.GameData.Teams.Home.Name
	linescore := lipgloss.NewStyle().Margin(1, 0, 0, 0).Render(m.buildLineScore())
	matchup := lipgloss.NewStyle().MarginTop(1).Render(m.renderMatchup())
	plays := lipgloss.NewStyle().MarginTop(1).Render(m.plays.View())

	content := lipgloss.JoinVertical(lipgloss.Center, linescore, matchup, plays)

	all := lipgloss.JoinVertical(lipgloss.Center, header, content)

	return headerStyle.Width(width).Render(all)
}

func (m Model) buildLineScore() string {
	inningRow := []string{data.RenderInningState(m.game.Bug.InningSt)}
	awayRow := []string{m.game.Bug.AwayAbbr}
	homeRow := []string{m.game.Bug.HomeAbbr}
	linescore := m.game.Bug.Feed.LiveData.Linescore

	// build linescore rows
	for _, v := range linescore.Innings {
		inningRow = append(inningRow, strconv.Itoa(v.Num))
		awayRow = append(awayRow, strconv.Itoa(v.Away.Runs))
		homeRow = append(homeRow, strconv.Itoa(v.Home.Runs))
	}
	// fill rows to show 9 innings minimum
	if len(inningRow) < 10 {
		for i := m.game.Bug.Inning + 1; i < 10; i++ {
			inningRow = append(inningRow, strconv.Itoa(i))
			awayRow = append(awayRow, "0")
			homeRow = append(homeRow, "0")
		}
	}

	// add runs, hits, errors cols and values
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
				switch col {
				case m.game.Bug.Inning:
					return linescoreStyle.Foreground(lipgloss.Yellow).Bold(true)
				}
				return linescoreStyle.Foreground(lipgloss.Magenta).Bold(true)
			}
			switch col {
			case 0:
				return linescoreStyle.Foreground(lipgloss.Magenta).Bold(true)
			case m.game.Bug.Inning:
				return linescoreStyle.Foreground(lipgloss.Yellow).Bold(true)
			}
			return linescoreStyle
		}).
		Headers(inningRow...).Rows(rows...)
	return t.Render()
}
