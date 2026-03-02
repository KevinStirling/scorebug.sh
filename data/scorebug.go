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

	Inning       int
	InningSt     string
	InningTop    string
	InningBottom string

	Out1 string
	Out2 string
	Out3 string

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

			Inning:       g.Linescore.CurrentInning,
			InningSt:     g.Linescore.InningState,
			InningTop:    setInningArrow(g.Linescore.InningState, "Top"),
			InningBottom: setInningArrow(g.Linescore.InningState, "Bottom"),

			Out1: SetOut(g.Linescore.Outs, 1),
			Out2: SetOut(g.Linescore.Outs, 2),
			Out3: SetOut(g.Linescore.Outs, 3),

			// Defaults
			On1B: "◇", On2B: "◇", On3B: "◇",
		}

		// enrich from live feed if present
		if s.Feed != nil {
			feed := *s.Feed

			bug.Balls = feed.LiveData.Plays.CurrentPlay.Count.Balls
			bug.Strikes = feed.LiveData.Plays.CurrentPlay.Count.Strikes

			setRunners(feed, &bug)
			setCurrentBP(feed, &bug)
		}
		out = append(out, bug)
	}
	return out
}

func playerKey(id int) string { return fmt.Sprintf("ID%d", id) }

func setRunners(f mlbstats.Feed, bug *ScoreBug) {
	offense := f.LiveData.Linescore.Offense

	bug.On1B = RunnerOn(offense.First != nil)
	bug.On2B = RunnerOn(offense.Second != nil)
	bug.On3B = RunnerOn(offense.Third != nil)
}

func setCurrentBP(f mlbstats.Feed, bug *ScoreBug) {
	b := f.LiveData.Plays.CurrentPlay.MatchUp.Batter
	p := f.LiveData.Plays.CurrentPlay.MatchUp.Pitcher

	splitName := func(name string) string {
		s := strings.Split(name, " ")
		if len(s) >= 2 {
			return s[1]
		}

		return name
	}

	bug.BatterName, bug.PitcherName = splitName(b.FullName), splitName(p.FullName)

	keyB, keyP := playerKey(b.Id), playerKey(p.Id)

	teams := []map[string]mlbstats.Player{f.LiveData.Boxscore.Teams.Home.Players, f.LiveData.Boxscore.Teams.Away.Players}

	for _, pm := range teams {
		if pl, ok := pm[keyB]; ok {
			if pl.SeasonStats.Batting.Avg != "" {
				bug.BatterAvg = pl.SeasonStats.Batting.Avg
				break
			}
		}
	}

	for _, pm := range teams {
		if pl, ok := pm[keyP]; ok {
			if pl.Stats.Pitching.PitchesThrown != 0 {
				bug.PitchCount = pl.Stats.Pitching.PitchesThrown
				break
			}
		}
	}
}
