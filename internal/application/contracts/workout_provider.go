package contracts

import (
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type WorkoutProvider interface {
	FetchWorkout(unixDate int64) (*domain.Workout, error)
}
