package mlbstats

// Feed is the struct representation of the json returned from an mlbstats api request to get a live game feed
type Feed struct {
	GameData struct {
		Teams struct {
			Home struct {
				Abbreviation string `json:"abbreviation"`
				Name         string `json:"name"`
			} `json:"home"`
			Away struct {
				Abbreviation string `json:"abbreviation"`
				Name         string `json:"name"`
			} `json:"away"`
		} `json:"teams"`
	} `json:"gameData"`

	LiveData struct {
		Linescore struct {
			CurrentInning int    `json:"currentInning"`
			InningState   string `json:"inningState"`
			Outs          int    `json:"outs"`
			Innings       []struct {
				Num        int    `json:"num"`
				OrdinalNum string `json:"ordinalNum"`
				Home       struct {
					Errors     int `json:"errors"`
					Hits       int `json:"hits"`
					LeftOnBase int `json:"leftOnBase"`
					Runs       int `json:"runs"`
				} `json:"home"`
				Away struct {
					Errors     int `json:"errors"`
					Hits       int `json:"hits"`
					LeftOnBase int `json:"leftOnBase"`
					Runs       int `json:"runs"`
				} `json:"away"`
			} `json:"innings"`
			Teams struct {
				Home struct {
					Runs   int `json:"runs"`
					Hits   int `json:"hits"`
					Errors int `json:"errors"`
				} `json:"home"`
				Away struct {
					Runs   int `json:"runs"`
					Hits   int `json:"hits"`
					Errors int `json:"errors"`
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
			AllPlays []struct {
				Type   string `json:"type"`
				Result struct {
					Description string `json:"description"`
					IsOut       bool   `json:"isOut"`
					HasReview   bool   `json:"hasReview"`
					EventType   string `json:"eventType"`
				} `json:"result"`
				About struct {
					Inning     int    `json:"inning"`
					HalfInning string `json:"halfInning"`
				}
			} `json:"allPlays"`
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
					Batside struct {
						Code        string `json:"code"`
						Description string `json:"description"`
					} `json:"batside"`
					PitchHand struct {
						Code        string `json:"code"`
						Description string `json:"description"`
					} `json:"pitchHand"`
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

func (f Feed) Player(id int) (Player, bool) {
	pKey := FormatPlayerKey(id)
	teams := []map[string]Player{
		f.LiveData.Boxscore.Teams.Home.Players,
		f.LiveData.Boxscore.Teams.Away.Players,
	}

	for _, players := range teams {
		if p, ok := players[pKey.String()]; ok {
			return p, true
		}
	}
	return Player{}, false
}
