package strava

import (
	"context"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/usecase/input"
)

type Client interface {
	GetActivitiesByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error)
	GetActivitiesByPeriod(ctx context.Context, period input.Period) ([]*ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
}
