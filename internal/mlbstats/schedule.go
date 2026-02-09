package mlbstats

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
