package contracts

import (
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type WorkoutProvider interface {
	GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error)
}
