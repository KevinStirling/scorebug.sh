package scorebug

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
)

// Renders the current batter / pitcher stat strings for a given ScoreBug.
// isHome determines which side should show batter vs pitcher depending on inning.
func renderBp(game data.ScoreBug, isHome bool) string {
	side := strings.ToLower(game.InningSt)

	if game.PitcherName == "" || game.BatterName == "" {
		if game.Status == "Final" && isHome {
			return "Final"
		}
		return ""
	}

	switch side {
	case "top":
		if isHome {
			return fmt.Sprintf("%s %dp", game.PitcherName, game.PitchCount)
		}
		return fmt.Sprintf("%s %s", game.BatterName, game.BatterAvg)

	case "bottom":
		if !isHome {
			return fmt.Sprintf("%s %dp", game.PitcherName, game.PitchCount)
		}
		return fmt.Sprintf("%s %s", game.BatterName, game.BatterAvg)

	default:
		return ""
	}
}

// Render renders a single scorebug from a flattened ScoreBug.
func Render(game data.ScoreBug) string {
	rows := [][]string{
		{
			game.HomeAbbr,
			game.AwayAbbr,
			game.On2B,
			data.SetOut(game.Outs, 1),
			func() string {
				if game.InningSt == "Top" {
					return game.InningArrow
				}
				return ""
			}(),
		},
		{
			strconv.Itoa(game.HomeRuns),
			strconv.Itoa(game.AwayRuns),
			game.On3B + " - " + game.On1B,
			data.SetOut(game.Outs, 2),
			strconv.Itoa(game.Inning),
		},
		{
			renderBp(game, true),
			renderBp(game, false),
			fmt.Sprintf("%s", strconv.Itoa(game.Balls)+"-"+strconv.Itoa(game.Strikes)),
			data.SetOut(game.Outs, 3),
			func() string {
				if game.InningSt == "Bottom" {
					return game.InningArrow
				}
				return ""
			}(),
		},
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

	return fmt.Sprintf("%s", t)
}
