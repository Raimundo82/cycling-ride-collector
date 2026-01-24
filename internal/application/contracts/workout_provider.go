package contracts

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type WorkoutProvider interface {
	GetWorkoutsByDate(date time.Time) ([]*domain.Workout, error)
}
