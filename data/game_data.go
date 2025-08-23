package data

import (
	"encoding/json"
	"net/http"
)

// TODO currently the Feed is only used to build the ScoreBug, however this struct
// can be expanded to parse more of the gamedata, and shared by other structs to for
// building data formatted for different UI formats
type Feed struct {
	GameData struct {
		Teams struct {
			Home struct {
				Abbreviation string `json:"abbreviation"`
			} `json:"home"`
			Away struct {
				Abbreviation string `json:"abbreviation"`
			} `json:"away"`
		} `json:"teams"`
	} `json:"gameData"`

	LiveData struct {
		Linescore struct {
			CurrentInning int    `json:"currentInning"`
			InningState   string `json:"inningState"`
			Outs          int    `json:"outs"`
			Teams         struct {
				Home struct {
					Runs int `json:"runs"`
				} `json:"home"`
				Away struct {
					Runs int `json:"runs"`
				} `json:"away"`
			} `json:"teams"`
			Offense struct {
				First *struct {
					ID int `json:"id"`
				} `json:"first"`
				Second *struct {
					ID int `json:"id"`
				} `json:"second"`
				Third *struct {
					ID int `json:"id"`
				} `json:"third"`
			} `json:"offense"`
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
	} `json:"liveData"`
}

type ScoreBug struct {
	HomeAbbr string
	AwayAbbr string
	HomeRuns int
	AwayRuns int
	Inning   int
	InningSt string
	Outs     int
	Balls    int
	Strikes  int
	On1B     string
	On2B     string
	On3B     string
}

func GetGameFeed() Feed {
	resp, _ := http.Get("https://statsapi.mlb.com/api/v1.1/game/776796/feed/live")
	defer resp.Body.Close()

	var feed Feed
	json.NewDecoder(resp.Body).Decode(&feed)

	return feed
}

func BuildScoreBug(f Feed) ScoreBug {
	setRunner := func(s []int, i int) string {
		if i >= 0 && i < len(s) && s[i] != 0 {
			return "●"
		}
		return " "
	}

	setInningState := func(state string) string {
		switch state {
		case "Top":
			return "↑"
		case "Bottom":
			return "↓"
		default:
			return ""
		}
	}

	return ScoreBug{
		HomeAbbr: f.GameData.Teams.Home.Abbreviation,
		AwayAbbr: f.GameData.Teams.Away.Abbreviation,
		HomeRuns: f.LiveData.Linescore.Teams.Home.Runs,
		AwayRuns: f.LiveData.Linescore.Teams.Away.Runs,
		Inning:   f.LiveData.Linescore.CurrentInning,
		InningSt: setInningState(f.LiveData.Linescore.InningState),
		Outs:     f.LiveData.Linescore.Outs,
		Balls:    f.LiveData.Plays.CurrentPlay.Count.Balls,
		Strikes:  f.LiveData.Plays.CurrentPlay.Count.Strikes,
		On1B:     setRunner(f.LiveData.Plays.CurrentPlay.RunnerIndex, 1),
		On2B:     setRunner(f.LiveData.Plays.CurrentPlay.RunnerIndex, 2),
		On3B:     setRunner(f.LiveData.Plays.CurrentPlay.RunnerIndex, 3),
	}
}
