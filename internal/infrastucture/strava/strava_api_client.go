package strava

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type StravaApiClient interface {
	GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
}
