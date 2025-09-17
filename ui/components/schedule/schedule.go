package schedule

import (
	"fmt"
	"strconv"
	"strings"
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

func renderBp(g data.Game, isHome bool) string {
	var side = strings.ToLower(g.InningSt)
	switch side {
	case "top":
		if isHome {
			return fmt.Sprintf("%s, P: %d", g.Pitcher, g.PitchCount)
		} else {
			return fmt.Sprintf("%s, %s", g.Batter, g.BatterAvg)
		}
	case "bottom":
		if !isHome {
			return fmt.Sprintf("%s, P: %d", g.Pitcher, g.PitchCount)
		} else {
			return fmt.Sprintf("%s, %s", g.Batter, g.BatterAvg)
		}
	default:
		return ""
	}
}

func renderSchedule(g data.Schedule) string {
	var s string
	if len(g.Games) > 0 {
		//TODO put this in a dedicated "render scorebug" func so it can be used by other commands
		for _, g := range g.Games {
			var rows [][]string
			rows = [][]string{
				// move this annoymous func out of here, not readable enough
				{g.HomeAbbr, g.AwayAbbr, g.On2B, data.SetOut(g.Outs, 1), func() string {
					if g.InningSt == "Top" {
						return g.InningArrow
					}
					return ""
				}()}, {strconv.Itoa(g.HomeRuns), strconv.Itoa(g.AwayRuns), g.On3B + " - " + g.On1B, data.SetOut(g.Outs, 2), strconv.Itoa(g.Inning)},
				{renderBp(g, true), renderBp(g, false),
					fmt.Sprintf("%s", strconv.Itoa(g.Balls)+"-"+strconv.Itoa(g.Strikes)), data.SetOut(g.Outs, 3),
					func() string {
						if g.InningSt == "Bottom" {
							return g.InningArrow
						}
						return ""
					}()},
			}
			var (
				purple    = lipgloss.Color("99")
				cellStyle = lipgloss.NewStyle().Padding(0, 1)
				outsCol   = lipgloss.NewStyle().BorderLeft(false).Width(3).Align(lipgloss.Center)
			)
			t := table.New().
				Width(58).
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(purple)).
				StyleFunc(func(row, col int) lipgloss.Style {
					switch col {
					case 2:
						return lipgloss.NewStyle().Width(9).Align(lipgloss.Center)
					case 3:
						return outsCol
					case 4:
						return cellStyle.Width(3)
					}
					return cellStyle.Align(lipgloss.Center)
				}).
				Rows(rows...)
			// TODO Use JoinHorizonal and JoinVertical in combination with putting the
			// scorebugs in slices (one for each game status) to format them easier
			s += fmt.Sprintf("%s\n", t)
		}
	}
	return s
}
