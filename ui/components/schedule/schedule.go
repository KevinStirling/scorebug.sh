package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
)

type Model struct {
	games     data.Schedule
	paginator paginator.Model
	err       error
}

func NewModel() Model {
	games := data.BuildSchedule(data.GetSchedule())
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(games.Games))
	return Model{
		paginator: p,
		games:     games,
	}
}

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
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case data.Schedule:
		m.games = msg
		return m, nil

	case tickMsg:
		return m, tea.Batch(checkServer(), tickAfter(10*time.Second))

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	g := renderSchedule(m.games)
	var b strings.Builder
	b.WriteString("\n Games \n")
	start, end := m.paginator.GetSliceBounds(len(g))
	for _, item := range g[start:end] {
		b.WriteString("\n" + item)
	}
	b.WriteString("\n" + m.paginator.View())
	b.WriteString("\n\n  h/l ←/→ page • q: quit\n")
	return b.String()
}

func renderBp(g data.Game, isHome bool) string {
	var side = strings.ToLower(g.InningSt)
	if len(g.Pitcher) == 0 || len(g.Batter) == 0 {
		if g.Status == "Final" && isHome {
			return "Final"
		}
		return ""
	}
	switch side {
	case "top":
		if isHome {
			return fmt.Sprintf("%s %dp", g.Pitcher, g.PitchCount)
		} else {
			return fmt.Sprintf("%s %s", g.Batter, g.BatterAvg)
		}
	case "bottom":
		if !isHome {
			return fmt.Sprintf("%s %dp", g.Pitcher, g.PitchCount)
		} else {
			return fmt.Sprintf("%s %s", g.Batter, g.BatterAvg)
		}
	default:
		return ""
	}
}

func renderSchedule(g data.Schedule) []string {
	var bugCells []string
	if len(g.Games) > 0 {
		// TODO move this logic into scorebug component, update scorebug component to accept more modular input
		for _, game := range g.Games {
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
				purple    = lipgloss.Color("99")
				cellStyle = lipgloss.NewStyle().Padding(0, 1)
				outsCol   = lipgloss.NewStyle().BorderLeft(false).Width(3).Align(lipgloss.Center)
			)
			t := table.New().
				Width(SB_WIDTH).
				Height(SB_HEIGHT).
				Border(lipgloss.NormalBorder()).
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
			bugCells = append(bugCells, bugStr)
		}
	}
	return bugCells
}
