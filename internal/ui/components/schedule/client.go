package schedule

import (
	"time"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
)

// ScheduleClient is an interface for retrieving schedule and game data from mlbstats
type ScheduleClient interface {
	// Schedule returns the schedule for the given date
	Schedule(date time.Time) (mlbstats.Schedule, error)
	// GameFeed returns the game feed for the given gameLink
	GameFeed(gameLink string) (mlbstats.Feed, error)
}
