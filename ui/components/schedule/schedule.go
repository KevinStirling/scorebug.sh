package schedule

import (
	"fmt"
	"strconv"
	"time"

	"github.com/KevinStirling/scorebug.sh/data"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Model struct {
	games data.Schedule
	err   error
}

func NewModel() Model { return Model{} }

type tickMsg time.Time

func tickAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func checkServer() tea.Cmd {
	return func() tea.Msg {
		return data.BuildSchedule(data.GetSchedule())
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(checkServer(), tickAfter(10*time.Second))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case data.Schedule:
		m.games = msg
		return m, nil

	case tickMsg:
		return m, tea.Batch(checkServer(), tickAfter(10*time.Second))

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return "\n" + renderSchedule(m.games) + "\n"
}

func renderSchedule(g data.Schedule) string {
	var s string
	if len(g.Games) > 0 {
		for _, g := range g.Games {
			rows := [][]string{
				{g.HomeAbbr, g.AwayAbbr, strconv.Itoa(g.Outs) + " OUTS", "[" + g.On2B + "]"},
				{strconv.Itoa(g.HomeRuns), strconv.Itoa(g.AwayRuns), strconv.Itoa(g.Balls) + "-" + strconv.Itoa(g.Strikes), "[" + g.On3B + "] _ [" + g.On1B + "]", g.InningSt + strconv.Itoa(g.Inning)},
			}
			var (
				purple    = lipgloss.Color("99")
				cellStyle = lipgloss.NewStyle().Padding(0, 1)
			)
			t := table.New().
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(purple)).
				StyleFunc(func(row, col int) lipgloss.Style { return cellStyle.Align(lipgloss.Center) }).
				Rows(rows...)
			s += fmt.Sprintf("%s\n", t)
		}
	}
	return s
}
