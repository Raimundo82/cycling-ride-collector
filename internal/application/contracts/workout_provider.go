package contracts

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/usecase/input"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type SingleWorkoutProvider interface {
	GetWorkoutsByDate(date time.Time) ([]*domain.Workout, error)
}

type PeriodWorkoutProvider interface {
	GetWorkoutsByPeriod(period input.Period) ([]*domain.Workout, error)
}

type WorkoutProvider interface {
	GetWorkoutsByDate(date time.Time) ([]*domain.Workout, error)
	GetWorkoutsByPeriod(period input.Period) ([]*domain.Workout, error)
}
