package contracts

import "github.com/raimundo82/go-strava-weekly/internal/domain"

type WorkoutProvider interface {
	FetchWorkouts() ([]*domain.Workout, error)
}
