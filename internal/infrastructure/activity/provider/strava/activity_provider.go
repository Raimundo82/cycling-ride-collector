package strava

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
)

type ActivityProvider interface {
	GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*model.ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*model.DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*model.WattsStreamDto, error)
}
