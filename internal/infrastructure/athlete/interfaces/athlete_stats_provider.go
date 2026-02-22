package athlete_interfaces

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
)

type AthleteStatsProvider interface {
	GetDetailedAthlete(ctx context.Context) (*model.DetailedAthlete, error)
	GetAthleteZones(ctx context.Context) (*model.Zones, error)
}
