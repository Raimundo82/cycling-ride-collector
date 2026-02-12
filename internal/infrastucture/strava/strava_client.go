package strava

import (
	"context"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type StravaClient interface {
	GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
}
