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
			rows := [][]string{
				// move this annoymous func out of here, not readable enough
				{g.HomeAbbr, g.AwayAbbr, "[" + g.On2B + "]", func() string {
					if g.InningSt == "Top" {
						return g.InningArrow
					}
					return ""
				}()}, {strconv.Itoa(g.HomeRuns), strconv.Itoa(g.AwayRuns), "[" + g.On3B + "] _ [" + g.On1B + "]", strconv.Itoa(g.Inning)},
				{renderBp(g, true), renderBp(g, false),
					fmt.Sprintf("%s   %s", strconv.Itoa(g.Balls)+"-"+strconv.Itoa(g.Strikes), strconv.Itoa(g.Outs)+" OUTS"),
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
			)
			t := table.New().
				Width(69).
				Border(lipgloss.NormalBorder()).
				BorderStyle(lipgloss.NewStyle().Foreground(purple)).
				StyleFunc(func(row, col int) lipgloss.Style {
					switch col {
					case 3:
						return cellStyle.Width(3)
					}
					return cellStyle.Align(lipgloss.Center)
				}).
				Rows(rows...)
			s += fmt.Sprintf("%s\n", t)
		}
	}
	return s
}
