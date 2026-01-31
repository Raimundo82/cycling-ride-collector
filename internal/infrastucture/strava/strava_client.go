package strava

import (
	"context"
	"time"
)

type Client interface {
	GetActivitiesByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error)
	GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
	GetAthleteData(ctx context.Context) (*AthleteDto, error)
	GetAthleteZones(ctx context.Context) (*AthleteZonesDto, error)
}
