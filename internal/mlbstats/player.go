package mlbstats

import "fmt"

// Player is struct representation of a player in the json returned from an mlbstats api request to get a live game feed
// In the json response for a live game request, this is found nested inside liveData.boxscore.teams.{away||home}.players
type Player struct {
	Person struct {
		Id       int    `json:"id"`
		FullName string `json:"fullName"`
	} `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     struct {
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
	Stats struct {
		Pitching struct {
			PitchesThrown int `json:"pitchesThrown"`
		} `json:"pitching"`
	} `json:"stats"`
	SeasonStats struct {
		Pitching struct {
			Win        int    `json:"wins"`
			Losses     int    `json:"losses"`
			Era        string `json:"era"`
			StrikeOuts int    `json:"strikeOuts"`
			Whip       string `json:"whip"`
		}
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

type PlayerKey struct {
	value string
}

func (p PlayerKey) String() string { return p.value }

func FormatPlayerKey(id int) PlayerKey {
	return PlayerKey{value: (fmt.Sprintf("ID%d", id))}
}
