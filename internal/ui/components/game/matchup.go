package game

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

type Pitcher struct {
	Name       string
	PitchHand  string
	Number     string
	WinLoss    string
	Era        string
	StrikeOuts string
	Whip       string
}

type Batter struct {
	Name     string
	BatSide  string
	Number   string
	Avg      string
	Ops      string
	Rbi      string
	HomeRuns string
	Position string
}

type MatchUp struct {
	Pitcher Pitcher
	Batter  Batter
}

func (m Model) renderMatchup() string {
	matchup := m.buildMatchup()
	var pitcherHeader string

	switch m.game.Bug.InningSt {
	case "Top":
		pitcherHeader = fmt.Sprintf("Pitching for %s", m.game.Bug.HomeAbbr)
	case "Bottom":
		pitcherHeader = fmt.Sprintf("Pitching for %s", m.game.Bug.AwayAbbr)
	}

	pitcherRows := [][]string{
		{matchup.Pitcher.Name, matchup.Pitcher.StrikeOuts, "K"},
		{matchup.Pitcher.PitchHand + "HP #" + matchup.Pitcher.Number, matchup.Pitcher.Era, "ERA"},
		{"", matchup.Pitcher.WinLoss, "W/L"},
		{"", matchup.Pitcher.Whip, "WHIP"},
	}

	pt := table.New().
		BorderColumn(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				switch row {
				case 0:
					return cellStyle.Foreground(lipgloss.Yellow)
				case 1:
					return cellStyle.Foreground(lipgloss.BrightBlack)
				}
			case 1:
				return cellStyle.Align(lipgloss.Right)
			}
			return cellStyle
		}).Rows(pitcherRows...)

	var batterHeader string

	switch m.game.Bug.InningSt {
	case "Top":
		batterHeader = fmt.Sprintf("Batting for %s", m.game.Bug.AwayAbbr)
	case "Bottom":
		batterHeader = fmt.Sprintf("Batting for %s", m.game.Bug.HomeAbbr)
	}

	batterRows := [][]string{
		{matchup.Batter.Avg, "AVG", matchup.Batter.Name},
		{matchup.Batter.Ops, "OPS", matchup.Batter.Position + " #" + matchup.Batter.Number + " (" + matchup.Batter.BatSide + ")"},
		{matchup.Batter.Rbi, "RBI"},
		{matchup.Batter.HomeRuns, "HR"},
	}

	bt := table.New().
		BorderColumn(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch col {
			case 0:
				return cellStyle.Align(lipgloss.Right)
			case 2:
				switch row {
				case 0:
					return cellStyle.Foreground(lipgloss.Yellow).Align(lipgloss.Right)
				case 1:
					return cellStyle.Foreground(lipgloss.BrightBlack).Align(lipgloss.Right)
				}
			}
			return cellStyle
		}).Rows(batterRows...)

	batterCard := lipgloss.JoinVertical(lipgloss.Right, lipgloss.NewStyle().Align(lipgloss.Right).PaddingRight(1).Foreground(lipgloss.Magenta).Render(batterHeader), bt.Render())

	pitcherCard := lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Align(lipgloss.Left).PaddingLeft(1).Foreground(lipgloss.Magenta).Render(pitcherHeader), pt.Render())

	return lipgloss.JoinHorizontal(lipgloss.Top, pitcherCard, batterCard)
}

func (m Model) buildMatchup() MatchUp {
	var pitcherStats Pitcher
	var batterStats Batter

	matchupData := m.game.Bug.Feed.LiveData.Plays.CurrentPlay.MatchUp

	pitcher, ok := m.game.Bug.Feed.Player(matchupData.Pitcher.Id)

	if ok {
		pitcherStats = Pitcher{
			Name:       pitcher.Person.FullName,
			PitchHand:  matchupData.PitchHand.Code,
			Number:     pitcher.JerseyNumber,
			WinLoss:    strconv.Itoa(pitcher.SeasonStats.Pitching.Win) + "-" + strconv.Itoa(pitcher.SeasonStats.Pitching.Losses),
			Era:        pitcher.SeasonStats.Pitching.Era,
			StrikeOuts: strconv.Itoa(pitcher.SeasonStats.Pitching.StrikeOuts),
			Whip:       pitcher.SeasonStats.Pitching.Whip,
		}
	}

	batter, ok := m.game.Bug.Feed.Player(matchupData.Batter.Id)
	if ok {
		batterStats = Batter{
			Name:     batter.Person.FullName,
			BatSide:  matchupData.Batside.Code,
			Number:   batter.JerseyNumber,
			Avg:      batter.SeasonStats.Batting.Avg,
			Ops:      batter.SeasonStats.Batting.Ops,
			Rbi:      strconv.Itoa(batter.SeasonStats.Batting.Rbi),
			HomeRuns: strconv.Itoa(batter.SeasonStats.Batting.HomeRuns),
			Position: batter.Position.Abbreviation,
		}
	}

	return MatchUp{Pitcher: pitcherStats, Batter: batterStats}
}
