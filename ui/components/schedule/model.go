package schedule

import (
	"time"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
)

type StatsClient interface {
	Schedule(date *time.Time) (mlbstats.Schedule, error)
	// LiveFeed(link string) (mlbstats.Feed, error)
}
