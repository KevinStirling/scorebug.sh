package schedule

import (
	"time"

	"github.com/KevinStirling/scorebug.sh/internal/mlbstats"
)

type ScheduleClient interface {
	Schedule(date *time.Time) (mlbstats.Schedule, error)
	GameFeed(gameLink string) (mlbstats.Feed, error)
}
