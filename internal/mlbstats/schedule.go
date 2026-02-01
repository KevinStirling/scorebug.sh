package mlbstats

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

// Games stores only the useful data pulled from the mlb statsapi
type Schedule struct {
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
				Offense       struct {
					First *struct {
						ID int
					} `json:"first"`
					Second *struct {
						ID int
					} `json:"second"`
					Third *struct {
						ID int
					} `json:"third"`
				}
				Teams struct {
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

// Retrieves a schedule of games for the given date. Each game on the schedule
// contains the data required to render a scorebug
// If data is not passed, time.Now will be used
func GetSchedule(date *time.Time) (Schedule, error) {
	var d time.Time
	if date == nil {
		d = time.Now()
	} else {
		d = *date
	}
	scheduleUrl = fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s&hydrate=linescore,team", d.Format("01/02/2006"))

	resp, err := http.Get(scheduleUrl)
	if err != nil {
		return Schedule{}, err
	}
	defer resp.Body.Close()

	var schedule Schedule
	json.NewDecoder(resp.Body).Decode(&schedule)

	return schedule, nil
}
