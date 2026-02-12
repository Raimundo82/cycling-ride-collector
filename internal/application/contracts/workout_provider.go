package contracts

import (
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type WorkoutProvider interface {
	GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error)
}
