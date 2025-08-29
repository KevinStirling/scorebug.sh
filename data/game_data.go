package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
				MatchUp struct {
					Batter struct {
						Id       int    `json:"id"`
						FullName string `json:"fullName"`
					} `json:"batter"`
					Pitcher struct {
						Id       int    `json:"id"`
						FullName string `json:"fullName"`
					} `json:"pitcher"`
				} `json:"matchUp"`
				RunnerIndex []int `json:"runnerIndex"`
			} `json:"currentPlay"`
		} `json:"plays"`
		Boxscore struct {
			Teams struct {
				Home struct {
					Players map[string]Player `json:"players"`
				} `json:"home"`
				Away struct {
					Players map[string]Player `json:"players"`
				} `json:"away"`
			} `json:"teams"`
		} `json:"boxscore"`
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

type Player struct {
	Person struct {
		Id       int    `json:"id"`
		FullName string `json:"fullName"`
	} `json:"person"`
	Stats struct {
		Pitching struct {
			PitchesThrown int `json:"pitchesThrown"`
		} `json:"pitching"`
	} `json:"stats"`
	SeasonStats struct {
		Batting struct {
			Avg string `json:"avg"`
		} `json:"batting"`
	} `json:"seasonStats"`
}

type CurrentBP struct {
	BatterName  string
	BatterAvg   string
	PitcherName string
	PitchCount  int
}

func playerKey(id int) string { return fmt.Sprintf("ID%d", id) }

func getCurrentBP(f Feed) CurrentBP {
	bp := CurrentBP{}
	b := f.LiveData.Plays.CurrentPlay.MatchUp.Batter
	p := f.LiveData.Plays.CurrentPlay.MatchUp.Pitcher

	bp.BatterName, bp.PitcherName = strings.Split(b.FullName, " ")[1], strings.Split(p.FullName, " ")[1]

	keyB, keyP := playerKey(b.Id), playerKey(p.Id)

	teams := []map[string]Player{f.LiveData.Boxscore.Teams.Home.Players, f.LiveData.Boxscore.Teams.Away.Players}

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

func GetGameFeed(url string) Feed {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var feed Feed
	json.NewDecoder(resp.Body).Decode(&feed)

	return feed
}

func BuildScoreBug(f Feed) ScoreBug {
	return ScoreBug{
		HomeAbbr: f.GameData.Teams.Home.Abbreviation,
		AwayAbbr: f.GameData.Teams.Away.Abbreviation,
		HomeRuns: f.LiveData.Linescore.Teams.Home.Runs,
		AwayRuns: f.LiveData.Linescore.Teams.Away.Runs,
		Inning:   f.LiveData.Linescore.CurrentInning,
		InningSt: setInningArrow(f.LiveData.Linescore.InningState),
		Outs:     f.LiveData.Linescore.Outs,
		Balls:    f.LiveData.Plays.CurrentPlay.Count.Balls,
		Strikes:  f.LiveData.Plays.CurrentPlay.Count.Strikes,
		On1B:     setRunnerState(f.LiveData.Plays.CurrentPlay.RunnerIndex, 1),
		On2B:     setRunnerState(f.LiveData.Plays.CurrentPlay.RunnerIndex, 2),
		On3B:     setRunnerState(f.LiveData.Plays.CurrentPlay.RunnerIndex, 3),
	}
}
