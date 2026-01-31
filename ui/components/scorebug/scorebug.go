package scorebug

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/KevinStirling/scorebug.sh/data"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
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
	var s string

	// TODO should update to render a single game's scorebug.
	// needs refactor, should adapt to use the Game type for managing the data,
	// like the schedule does.

	// scorebug := Render(m.scorebug)
	// s = fmt.Sprintf("%s", scorebug)

	return "\n" + s + "\n"
}

// Renders the current batter / pitcher stat strings for a given Game type
// isHome is used to determine which team from `game` should be rendered.
// The `game.InningSt` is used to determine if a batter or pitcher should
// be rendered
func renderBp(game data.Game, isHome bool) string {
	var side = strings.ToLower(game.InningSt)
	if len(game.Pitcher) == 0 || len(game.Batter) == 0 {
		if game.Status == "Final" && isHome {
			return "Final"
		}
		return ""
	}
	switch side {
	case "top":
		if isHome {
			return fmt.Sprintf("%s %dp", game.Pitcher, game.PitchCount)
		} else {
			return fmt.Sprintf("%s %s", game.Batter, game.BatterAvg)
		}
	case "bottom":
		if !isHome {
			return fmt.Sprintf("%s %dp", game.Pitcher, game.PitchCount)
		} else {
			return fmt.Sprintf("%s %s", game.Batter, game.BatterAvg)
		}
	default:
		return ""
	}
}

// Renders a scorebug string for a give Game type
func Render(game data.Game) string {
	var rows [][]string
	rows = [][]string{
		{game.HomeAbbr, game.AwayAbbr, game.On2B, data.SetOut(game.Outs, 1), func() string {
			if game.InningSt == "Top" {
				return game.InningArrow
			}
			return ""
		}()},
		{strconv.Itoa(game.HomeRuns), strconv.Itoa(game.AwayRuns), game.On3B + " - " + game.On1B, data.SetOut(game.Outs, 2), strconv.Itoa(game.Inning)},
		{renderBp(game, true), renderBp(game, false),
			fmt.Sprintf("%s", strconv.Itoa(game.Balls)+"-"+strconv.Itoa(game.Strikes)), data.SetOut(game.Outs, 3),
			func() string {
				if game.InningSt == "Bottom" {
					return game.InningArrow
				}
				return ""
			}()},
	}
	var (
		purple    = lipgloss.Color("65")
		cellStyle = lipgloss.NewStyle().Padding(0, 1)
		outsCol   = lipgloss.NewStyle().BorderLeft(false).Width(3).Align(lipgloss.Center)
	)
	t := table.New().
		Width(SB_WIDTH).
		Height(SB_HEIGHT).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 2:
				return lipgloss.NewStyle().Width(9).Align(lipgloss.Center)
			case 3:
				return outsCol
			case 4:
				return lipgloss.NewStyle().Width(3).Align(lipgloss.Center)
			}
			return cellStyle.Align(lipgloss.Center)
		}).
		Rows(rows...)
	bugStr := fmt.Sprintf("%s", t)
	return bugStr
}
