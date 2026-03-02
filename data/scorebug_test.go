package data_test

import (
	"testing"

	"github.com/KevinStirling/scorebug.sh/data"
	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
	"github.com/KevinStirling/scorebug.sh/internal/snapshots"
)

func TestBuildScoreBugs_PreviewGame(t *testing.T) {
	snap := snapshots.GameSnapshot{
		Game: mlbstats.Game{
			Link: "/api/v1/game/12345/feed/live",
			Status: struct {
				AbstractGameState string `json:"abstractGameState"`
			}{AbstractGameState: "Preview"},
			Teams: makeTeams("NYY", "BOS"),
		},
		Feed: nil, // no feed for preview games
	}

	bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})

	if len(bugs) != 1 {
		t.Fatalf("expected 1 bug, got %d", len(bugs))
	}

	bug := bugs[0]

	// Check basic fields
	if bug.Status != "Preview" {
		t.Errorf("Status = %q, want %q", bug.Status, "Preview")
	}
	if bug.HomeAbbr != "NYY" {
		t.Errorf("HomeAbbr = %q, want %q", bug.HomeAbbr, "NYY")
	}
	if bug.AwayAbbr != "BOS" {
		t.Errorf("AwayAbbr = %q, want %q", bug.AwayAbbr, "BOS")
	}

	// Preview games should have empty bases (default diamonds)
	if bug.On1B != "◇" || bug.On2B != "◇" || bug.On3B != "◇" {
		t.Errorf("bases should be empty for preview, got On1B=%s On2B=%s On3B=%s",
			bug.On1B, bug.On2B, bug.On3B)
	}
}

func TestBuildScoreBugs_LiveGame_RunnersOnBase(t *testing.T) {
	cases := []struct {
		name                   string
		first, second, third   bool
		want1B, want2B, want3B string
	}{
		{"no runners", false, false, false, "◇", "◇", "◇"},
		{"runner on first", true, false, false, "◆", "◇", "◇"},
		{"runner on second", false, true, false, "◇", "◆", "◇"},
		{"runner on third", false, false, true, "◇", "◇", "◆"},
		{"runners on corners", true, false, true, "◆", "◇", "◆"},
		{"bases loaded", true, true, true, "◆", "◆", "◆"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feed := makeFeed(tc.first, tc.second, tc.third)
			snap := snapshots.GameSnapshot{
				Game: mlbstats.Game{
					Link: "/api/v1/game/12345/feed/live",
					Status: struct {
						AbstractGameState string `json:"abstractGameState"`
					}{AbstractGameState: "Live"},
					Teams: makeTeams("LAD", "SF"),
				},
				Feed: &feed,
			}

			bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})
			bug := bugs[0]

			if bug.On1B != tc.want1B {
				t.Errorf("On1B = %s, want %s", bug.On1B, tc.want1B)
			}
			if bug.On2B != tc.want2B {
				t.Errorf("On2B = %s, want %s", bug.On2B, tc.want2B)
			}
			if bug.On3B != tc.want3B {
				t.Errorf("On3B = %s, want %s", bug.On3B, tc.want3B)
			}
		})
	}
}

func TestBuildScoreBugs_LiveGame_CountAndOuts(t *testing.T) {
	feed := mlbstats.Feed{}
	feed.LiveData.Plays.CurrentPlay.Count.Balls = 3
	feed.LiveData.Plays.CurrentPlay.Count.Strikes = 2

	snap := snapshots.GameSnapshot{
		Game: mlbstats.Game{
			Link: "/api/v1/game/12345/feed/live",
			Status: struct {
				AbstractGameState string `json:"abstractGameState"`
			}{AbstractGameState: "Live"},
			Teams: makeTeams("CHC", "STL"),
			Linescore: struct {
				CurrentInning int    `json:"currentInning"`
				InningState   string `json:"inningState"`
				Balls         int    `json:"balls"`
				Strikes       int    `json:"strikes"`
				Outs          int    `json:"outs"`
				Offense       struct {
					First  *struct{ ID int } `json:"first"`
					Second *struct{ ID int } `json:"second"`
					Third  *struct{ ID int } `json:"third"`
				}
				Teams struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				} `json:"teams"`
			}{
				CurrentInning: 7,
				InningState:   "Top",
				Outs:          2,
			},
		},
		Feed: &feed,
	}

	bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})
	bug := bugs[0]

	if bug.Balls != 3 {
		t.Errorf("Balls = %d, want 3", bug.Balls)
	}
	if bug.Strikes != 2 {
		t.Errorf("Strikes = %d, want 2", bug.Strikes)
	}
	if bug.Inning != 7 {
		t.Errorf("Inning = %d, want 7", bug.Inning)
	}
	if bug.InningSt != "Top" {
		t.Errorf("InningSt = %q, want %q", bug.InningSt, "Top")
	}
	if bug.InningTop != "↑" {
		t.Errorf("InningTop = %q, want %q", bug.InningTop, "↑")
	}
	if bug.InningBottom != "" {
		t.Errorf("InningBottom = %q, want %q", bug.InningBottom, "")
	}
	// 2 outs: Out1 and Out2 should be filled, Out3 empty
	if bug.Out1 != "◉" {
		t.Errorf("Out1 = %q, want %q", bug.Out1, "◉")
	}
	if bug.Out2 != "◉" {
		t.Errorf("Out2 = %q, want %q", bug.Out2, "◉")
	}
	if bug.Out3 != "◯" {
		t.Errorf("Out3 = %q, want %q", bug.Out3, "◯")
	}
}

func TestBuildScoreBugs_LiveGame_Score(t *testing.T) {
	feed := mlbstats.Feed{}

	snap := snapshots.GameSnapshot{
		Game: mlbstats.Game{
			Link: "/api/v1/game/12345/feed/live",
			Status: struct {
				AbstractGameState string `json:"abstractGameState"`
			}{AbstractGameState: "Live"},
			Teams: makeTeams("HOU", "TEX"),
			Linescore: struct {
				CurrentInning int    `json:"currentInning"`
				InningState   string `json:"inningState"`
				Balls         int    `json:"balls"`
				Strikes       int    `json:"strikes"`
				Outs          int    `json:"outs"`
				Offense       struct {
					First  *struct{ ID int } `json:"first"`
					Second *struct{ ID int } `json:"second"`
					Third  *struct{ ID int } `json:"third"`
				}
				Teams struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				} `json:"teams"`
			}{
				CurrentInning: 5,
				InningState:   "Bottom",
				Teams: struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				}{
					Home: struct {
						Runs int `json:"runs"`
					}{Runs: 4},
					Away: struct {
						Runs int `json:"runs"`
					}{Runs: 2},
				},
			},
		},
		Feed: &feed,
	}

	bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})
	bug := bugs[0]

	if bug.HomeRuns != 4 {
		t.Errorf("HomeRuns = %d, want 4", bug.HomeRuns)
	}
	if bug.AwayRuns != 2 {
		t.Errorf("AwayRuns = %d, want 2", bug.AwayRuns)
	}
	if bug.InningTop != "" {
		t.Errorf("InningTop = %q, want %q", bug.InningTop, "")
	}
	if bug.InningBottom != "↓" {
		t.Errorf("InningBottom = %q, want %q", bug.InningBottom, "↓")
	}
}

func TestBuildScoreBugs_Outs(t *testing.T) {
	cases := []struct {
		name                         string
		outs                         int
		wantOut1, wantOut2, wantOut3 string
	}{
		{"0 outs", 0, "◯", "◯", "◯"},
		{"1 out", 1, "◉", "◯", "◯"},
		{"2 outs", 2, "◉", "◉", "◯"},
		{"3 outs", 3, "◉", "◉", "◉"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snap := snapshots.GameSnapshot{
				Game: mlbstats.Game{
					Link: "/api/v1/game/12345/feed/live",
					Status: struct {
						AbstractGameState string `json:"abstractGameState"`
					}{AbstractGameState: "Live"},
					Teams: makeTeams("NYY", "BOS"),
					Linescore: struct {
						CurrentInning int    `json:"currentInning"`
						InningState   string `json:"inningState"`
						Balls         int    `json:"balls"`
						Strikes       int    `json:"strikes"`
						Outs          int    `json:"outs"`
						Offense       struct {
							First  *struct{ ID int } `json:"first"`
							Second *struct{ ID int } `json:"second"`
							Third  *struct{ ID int } `json:"third"`
						}
						Teams struct {
							Home struct {
								Runs int `json:"runs"`
							} `json:"home"`
							Away struct {
								Runs int `json:"runs"`
							} `json:"away"`
						} `json:"teams"`
					}{
						Outs: tc.outs,
					},
				},
				Feed: nil,
			}

			bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})
			bug := bugs[0]

			if bug.Out1 != tc.wantOut1 {
				t.Errorf("Out1 = %q, want %q", bug.Out1, tc.wantOut1)
			}
			if bug.Out2 != tc.wantOut2 {
				t.Errorf("Out2 = %q, want %q", bug.Out2, tc.wantOut2)
			}
			if bug.Out3 != tc.wantOut3 {
				t.Errorf("Out3 = %q, want %q", bug.Out3, tc.wantOut3)
			}
		})
	}
}

func TestBuildScoreBugs_FinalGame(t *testing.T) {
	snap := snapshots.GameSnapshot{
		Game: mlbstats.Game{
			Link: "/api/v1/game/12345/feed/live",
			Status: struct {
				AbstractGameState string `json:"abstractGameState"`
			}{AbstractGameState: "Final"},
			Teams: makeTeams("ATL", "PHI"),
			Linescore: struct {
				CurrentInning int    `json:"currentInning"`
				InningState   string `json:"inningState"`
				Balls         int    `json:"balls"`
				Strikes       int    `json:"strikes"`
				Outs          int    `json:"outs"`
				Offense       struct {
					First  *struct{ ID int } `json:"first"`
					Second *struct{ ID int } `json:"second"`
					Third  *struct{ ID int } `json:"third"`
				}
				Teams struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				} `json:"teams"`
			}{
				CurrentInning: 9,
				Teams: struct {
					Home struct {
						Runs int `json:"runs"`
					} `json:"home"`
					Away struct {
						Runs int `json:"runs"`
					} `json:"away"`
				}{
					Home: struct {
						Runs int `json:"runs"`
					}{Runs: 5},
					Away: struct {
						Runs int `json:"runs"`
					}{Runs: 3},
				},
			},
		},
		Feed: nil, // no feed for final games
	}

	bugs := data.BuildScoreBugs([]snapshots.GameSnapshot{snap})
	bug := bugs[0]

	if bug.Status != "Final" {
		t.Errorf("Status = %q, want %q", bug.Status, "Final")
	}
	if bug.HomeRuns != 5 {
		t.Errorf("HomeRuns = %d, want 5", bug.HomeRuns)
	}
	if bug.AwayRuns != 3 {
		t.Errorf("AwayRuns = %d, want 3", bug.AwayRuns)
	}
	if bug.Inning != 9 {
		t.Errorf("Inning = %d, want 9", bug.Inning)
	}
}

// Helper functions to reduce test boilerplate

func makeTeams(home, away string) struct {
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
} {
	return struct {
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
	}{
		Home: struct {
			Team struct {
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
		}{
			Team: struct {
				Abbreviation string `json:"abbreviation"`
			}{Abbreviation: home},
		},
		Away: struct {
			Team struct {
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
		}{
			Team: struct {
				Abbreviation string `json:"abbreviation"`
			}{Abbreviation: away},
		},
	}
}

func makeFeed(first, second, third bool) mlbstats.Feed {
	var feed mlbstats.Feed

	if first {
		feed.LiveData.Linescore.Offense.First = &struct {
			ID int `json:"id"`
		}{ID: 123}
	}
	if second {
		feed.LiveData.Linescore.Offense.Second = &struct {
			ID int `json:"id"`
		}{ID: 456}
	}
	if third {
		feed.LiveData.Linescore.Offense.Third = &struct {
			ID int `json:"id"`
		}{ID: 789}
	}

	return feed
}
