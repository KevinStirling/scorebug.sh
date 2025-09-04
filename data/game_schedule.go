package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	scheduleUrl = fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s&hydrate=linescore,team", time.Now().Format("01/02/2006"))
	statsUrl    = "https://statsapi.mlb.com"
)

// TodaysGames will be used to store the result of the api call https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=08/24/2025&hydrate=linescore,team
type TodaysGames struct {
	Dates []struct {
		Games []struct {
			Link   string `json:"link"`
			Status struct {
				AbstractGameState string `json:"abstractGameState"`
			} `json:"status"`
			Teams struct {
				Home struct {
					Team struct {
						Abbreviation string `json:"abbreviation"`
					} `json:"team"`
				} `json:"home"`
				Away struct {
					Team struct {
						Abbreviation string `json:"abbreviation"`
					} `json:"team"`
				} `json:"away"`
			} `json:"teams"`
			Linescore struct {
				CurrentInning int    `json:"currentInning"`
				InningState   string `json:"inningState"`
				Balls         int    `json:"balls"`
				Strikes       int    `json:"strikes"`
				Outs          int    `json:"outs"`
				Teams         struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				} `json:"teams"`
			} `json:"linescore"`
			Plays struct {
				CurrentPlay struct {
					Count struct {
						Balls   int `json:"balls"`
						Strikes int `json:"strikes"`
						Outs    int `json:"outs"`
					} `json:"count"`
					RunnerIndex []int `json:"runnerIndex"`
				} `json:"currentPlay"`
			} `json:"plays"`
		} `json:"games"`
	} `json:"dates"`
}

type Game struct {
	Link        string
	Status      string
	HomeAbbr    string
	AwayAbbr    string
	Batter      string
	Pitcher     string
	BatterAvg   string
	PitchCount  int
	HomeRuns    int
	AwayRuns    int
	Inning      int
	InningArrow string
	InningSt    string
	Outs        int
	Balls       int
	Strikes     int
	On1B        string
	On2B        string
	On3B        string
}
type Schedule struct {
	Games []Game
}

func GetSchedule() TodaysGames {
	resp, _ := http.Get(scheduleUrl)
	defer resp.Body.Close()

	var schedule TodaysGames
	json.NewDecoder(resp.Body).Decode(&schedule)

	return schedule
}

func BuildSchedule(t TodaysGames) Schedule {
	var s Schedule
	for _, d := range t.Dates {
		for _, g := range d.Games {
			row := struct {
				Link        string
				Status      string
				HomeAbbr    string
				AwayAbbr    string
				Batter      string
				Pitcher     string
				BatterAvg   string
				PitchCount  int
				HomeRuns    int
				AwayRuns    int
				Inning      int
				InningArrow string
				InningSt    string
				Outs        int
				Balls       int
				Strikes     int
				On1B        string
				On2B        string
				On3B        string
			}{
				Link:        g.Link,
				Status:      g.Status.AbstractGameState,
				HomeAbbr:    g.Teams.Home.Team.Abbreviation,
				AwayAbbr:    g.Teams.Away.Team.Abbreviation,
				HomeRuns:    g.Linescore.Teams.Home.Runs,
				AwayRuns:    g.Linescore.Teams.Away.Runs,
				Inning:      g.Linescore.CurrentInning,
				InningArrow: setInningArrow(g.Linescore.InningState),
				InningSt:    g.Linescore.InningState,
				Outs:        g.Linescore.Outs,
				Balls:       g.Linescore.Balls,
				Strikes:     g.Linescore.Strikes,
				On1B:        " ",
				On2B:        " ",
				On3B:        " ",
			}
			if row.Status != "Live" {
				feed := GetGameFeed(statsUrl + row.Link)
				bp := getCurrentBP(feed)
				row.On1B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 1)
				row.On2B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 2)
				row.On3B = SetRunnerState(feed.LiveData.Plays.CurrentPlay.RunnerIndex, 3)
				row.Batter = bp.BatterName
				row.Pitcher = bp.PitcherName
				row.BatterAvg = bp.BatterAvg
				row.PitchCount = bp.PitchCount
			}
			s.Games = append(s.Games, row)
		}
	}

	return s
}
