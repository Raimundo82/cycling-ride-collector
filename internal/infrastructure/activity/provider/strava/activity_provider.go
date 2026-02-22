package activity_strava

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
)

type ActivityProvider interface {
	GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*activity_model.ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*activity_model.DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*activity_model.WattsStreamDto, error)
}
