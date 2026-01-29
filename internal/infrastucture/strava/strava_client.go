package strava

import (
	"context"
	"time"
)

type Client interface {
	GetActivitiesByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
}
