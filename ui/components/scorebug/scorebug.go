package scorebug

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/KevinStirling/scorebug.sh/data"
)

const (
	SB_WIDTH  = 58
	SB_HEIGHT = 5
)

// Render a scorebug from the data.ScoreBug
func Render(game data.ScoreBug) string {
	rows := [][]string{
		{
			game.HomeAbbr,
			game.AwayAbbr,
			game.On2B,
			game.Out1,
			game.InningTop,
		},
		{
			strconv.Itoa(game.HomeRuns),
			strconv.Itoa(game.AwayRuns),
			game.On3B + " - " + game.On1B,
			game.Out2,
			strconv.Itoa(game.Inning),
		},
		{
			renderBp(game, true),
			renderBp(game, false),
			fmt.Sprintf("%d-%d", game.Balls, game.Strikes),
			game.Out3,
			game.InningBottom,
		},
	}

	t := table.New().
		Width(SB_WIDTH).
		Height(SB_HEIGHT).
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Green)).
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

	return fmt.Sprint(t)
}

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
