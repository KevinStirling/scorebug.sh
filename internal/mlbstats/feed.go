package mlbstats

// Player is struct representation of a player in the json returned from an mlbstats api request to get a live game feed
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

// Feed is the struct representation of the json returned from an mlbstats api request to get a live game feed
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
