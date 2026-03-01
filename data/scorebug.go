package data

import (
	"fmt"
	"strings"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/internal/snapshots"
)

type ScoreBug struct {
	Link   string
	Status string

	HomeAbbr string
	AwayAbbr string

	HomeRuns int
	AwayRuns int

	Inning      int
	InningSt    string
	InningArrow string

	Outs    int
	Balls   int
	Strikes int

	On1B string
	On2B string
	On3B string

	BatterName  string
	BatterAvg   string
	PitcherName string
	PitchCount  int
}

type CurrentBP struct {
	BatterName  string
	BatterAvg   string
	PitcherName string
	PitchCount  int
}

func BuildScoreBugs(snaps []snapshots.GameSnapshot) []ScoreBug {
	out := make([]ScoreBug, 0, len(snaps))

	for _, s := range snaps {
		g := s.Game

		bug := ScoreBug{
			Link:   g.Link,
			Status: g.Status.AbstractGameState,

			HomeAbbr: g.Teams.Home.Team.Abbreviation,
			AwayAbbr: g.Teams.Away.Team.Abbreviation,

			HomeRuns: g.Linescore.Teams.Home.Runs,
			AwayRuns: g.Linescore.Teams.Away.Runs,

			Inning:      g.Linescore.CurrentInning,
			InningSt:    g.Linescore.InningState,
			InningArrow: setInningArrow(g.Linescore.InningState),

			Outs: g.Linescore.Outs,

			// Defaults
			On1B: "◇", On2B: "◇", On3B: "◇",
		}

		// enrich from live feed if present
		if s.Feed != nil {
			feed := *s.Feed

			bug.Balls = feed.LiveData.Plays.CurrentPlay.Count.Balls
			bug.Strikes = feed.LiveData.Plays.CurrentPlay.Count.Strikes

			// might want to update this func to accept a pointer and just change all the bases at once
			bug.On1B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 1)
			bug.On2B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 2)
			bug.On3B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 3)

			// same for this one...
			bp := getCurrentBP(feed)
			bug.BatterName = bp.BatterName
			bug.BatterAvg = bp.BatterAvg
			bug.PitcherName = bp.PitcherName
			bug.PitchCount = bp.PitchCount
		}
		out = append(out, bug)
	}
	return out
}

func playerKey(id int) string { return fmt.Sprintf("ID%d", id) }

// refactor this to take a pointer to a feed and modify the feed directly?
func getCurrentBP(f mlbstats.Feed) CurrentBP {
	bp := CurrentBP{}
	b := f.LiveData.Plays.CurrentPlay.MatchUp.Batter
	p := f.LiveData.Plays.CurrentPlay.MatchUp.Pitcher

	splitName := func(name string) string {
		s := strings.Split(name, " ")
		if len(s) >= 2 {
			return s[1]
		}

		return name
	}

	bp.BatterName, bp.PitcherName = splitName(b.FullName), splitName(p.FullName)

	keyB, keyP := playerKey(b.Id), playerKey(p.Id)

	teams := []map[string]mlbstats.Player{f.LiveData.Boxscore.Teams.Home.Players, f.LiveData.Boxscore.Teams.Away.Players}

	for _, pm := range teams {
		if pl, ok := pm[keyB]; ok {
			if pl.SeasonStats.Batting.Avg != "" {
				bp.BatterAvg = pl.SeasonStats.Batting.Avg
				break
			}
		}
	}

	for _, pm := range teams {
		if pl, ok := pm[keyP]; ok {
			if pl.Stats.Pitching.PitchesThrown != 0 {
				bp.PitchCount = pl.Stats.Pitching.PitchesThrown
				break
			}
		}
	}

	return bp

}
