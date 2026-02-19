package athlete_interfaces

import "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"

type AthleteStatsProvider interface {
	GetDetailedAthlete() (*model.DetailedAthlete, error)
	GetAthleteZones() (*model.Zones, error)
}
