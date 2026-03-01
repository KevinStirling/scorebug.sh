package snapshots

import (
	"log"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
)

type GameSnapshot struct {
	Game mlbstats.Game
	Feed *mlbstats.Feed
}

func Build(client interface {
	GameFeed(string) (mlbstats.Feed, error)
}, games mlbstats.Schedule) ([]GameSnapshot, error) {
	var out []GameSnapshot

	for _, d := range games.Dates {
		for _, g := range d.Games {
			snap := GameSnapshot{Game: g}

			if g.Status.AbstractGameState == "Live" {
				feed, err := client.GameFeed(g.Link)
				if err != nil {
					log.Default().Print("ERROR failed to retrieve game feed")
					return nil, err
				}
				snap.Feed = &feed
			}
			out = append(out, snap)
		}
	}
	return out, nil
}
