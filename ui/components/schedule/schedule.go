package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/ui/screen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
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
	return renderSchedule(m.games)
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
	// determine the terminal dimensions, add 1 to width & height since table lines are not included in the render
	gCols, gRows, _, err := screen.CalculateScorebugs(SB_WIDTH+1, SB_HEIGHT+1)
	if err != nil {
		fmt.Printf("Error setting scorebug grid size")
	}
	var bugCells [][]string
	if len(g.Games) > 0 {
		for _, game := range g.Games {
			if game.Status == "Live" {
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
				lines := strings.Split(bugStr, "\n")
				if len(lines) > 0 && lines[len(lines)-1] == "" {
					lines = lines[:len(lines)-1]
				}
				bugCells = append(bugCells, lines)
			}
		}
	}
	if len(bugCells) == 0 {
		return "no live games :("
	}
	var s string
	totalCells := gCols * gRows
	blankBug := make([]string, SB_HEIGHT)
	blankLine := strings.Repeat(" ", SB_WIDTH)
	for i := 0; i < SB_HEIGHT; i++ {
		blankBug[i] = blankLine
	}
	for len(bugCells) < totalCells {
		bugCells = append(bugCells, blankBug)
	}
	for row := 0; row < gRows; row++ {
		for lineIdx := 0; lineIdx < SB_HEIGHT; lineIdx++ {
			var rowLine string
			for col := 0; col < gCols; col++ {
				cellIndex := row*gCols + col
				rowLine += bugCells[cellIndex][lineIdx]
				if col < gCols-1 {
					rowLine += " "
				}
			}
			s += rowLine + "\n"
		}
	}
	return s
}
