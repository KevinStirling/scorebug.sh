package mlbstats

// People is struct representation of the response from mlbstats's people api
type People struct {
	Id      uint32 `json:"id"`
	BatSide struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"batSide"`
	FirstLastName string  `json:"firstLastName"`
	Stats         []Stats `json:"stats"`
}

type Stats struct {
	Splits []struct {
		Split struct {
			Stat struct {
				// batting
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

				// fielding

				// pitching

			} `json:"stat"`
		}
	} `json:"splits"`
	Type struct {
		DisplayName string `json:"displayName"`
	} `json:"type"`
	Group struct {
		DisplayName string `json:"displayName"`
	} `json:"group"`
}
