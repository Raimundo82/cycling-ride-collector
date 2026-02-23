package contracts

import "github.com/raimundo82/cycling-ride-collector/internal/domain"

type AthleteDataProvider interface {
	GetAthleteData() (*domain.Athlete, error)
}
