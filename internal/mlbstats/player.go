package mlbstats

// Player is struct representation of a player in the json returned from an mlbstats api request to get a live game feed
// In the json response for a live game request, this is found nested inside liveData.boxscore.teams.{away||home}.players
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
			AirOuts     int    `json:"airOuts"`
			AtBats      int    `json:"atBats"`
			Avg         string `json:"avg"`
			Ops         string `json:"ops"`
			Obs         string `json:"obs"`
			Hits        int    `json:"hits"`
			Rbi         int    `json:"rbi"`
			Runs        int    `json:"runs"`
			HomeRuns    int    `json:"homeRunes"`
			Slg         string `json:"slg"`
			StrikeOuts  int    `json:"strikeOuts"`
			StolenBases int    `json:"stolenBases"`
			TotalBases  int    `json:"totalBases"`
		} `json:"batting"`
	} `json:"seasonStats"`
}
