package scorebug

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/KevinStirling/scorebug.sh/data"
)

var (
	//TODO hard coded game feed for testing render, remove if this package is used for final renders
	statsUrl   = "https://statsapi.mlb.com"
	gamePkLink = "https://statsapi.mlb.com/api/v1.1/game/776796/feed/live"
)

type Model struct {
	scoreBug data.ScoreBug
	error    error
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return checkServer
}

func checkServer() tea.Msg {
	return data.BuildScoreBug(data.GetGameFeed(statsUrl + gamePkLink))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case data.ScoreBug:
		m.scoreBug = data.ScoreBug(msg)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil

}

func (m Model) View() string {
	s := fmt.Sprintf("Checking for box score...\n")
	if &m.scoreBug != nil {

		rows := [][]string{
			{m.scoreBug.HomeAbbr, m.scoreBug.AwayAbbr, strconv.Itoa(m.scoreBug.Outs) + " OUTS", "[" + m.scoreBug.On2B + "]"},
			{strconv.Itoa(m.scoreBug.HomeRuns), strconv.Itoa(m.scoreBug.AwayRuns), strconv.Itoa(m.scoreBug.Balls) + "-" + strconv.Itoa(m.scoreBug.Strikes), "[" + m.scoreBug.On3B + "] _ [" + m.scoreBug.On1B + "]", m.scoreBug.InningSt + strconv.Itoa(m.scoreBug.Inning)},
		}
		var (
			purple = lipgloss.Color("99")
			// gray   = lipgloss.Color("245")
			cellStyle = lipgloss.NewStyle().Padding(0, 1)
		)
		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(purple)).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch col {
				case 0:
					return cellStyle.Align(lipgloss.Center)
				case 1:
					return cellStyle.Align(lipgloss.Center)
				case 2:
					return cellStyle.Align(lipgloss.Center)
				case 3:
					return cellStyle.Align(lipgloss.Center)
				case 4:
					return cellStyle.Align(lipgloss.Center)
				}
				return cellStyle
			}).
			Rows(rows...)

		s = fmt.Sprintf("%s", t)
	}

	return "\n" + s + "\n"
}

// func main() {
// 	if _, err := tea.NewProgram(Model{}).Run(); err != nil {
// 		fmt.Printf("oy, ya cooked, mate - %s", err.Error())
// 		os.Exit(1)
// 	}
// }
