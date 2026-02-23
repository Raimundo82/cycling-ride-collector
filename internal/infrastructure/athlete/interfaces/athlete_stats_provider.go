package athlete_interfaces

import (
	"context"

	athlete_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
)

type AthleteStatsProvider interface {
	GetDetailedAthlete(ctx context.Context) (*athlete_model.DetailedAthlete, error)
	GetAthleteZones(ctx context.Context) (*athlete_model.Zones, error)
}
