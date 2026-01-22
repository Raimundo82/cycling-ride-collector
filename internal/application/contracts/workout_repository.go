package contracts

import "github.com/raimundo82/go-strava-weekly/internal/domain"

type WorkoutRepository interface {
	Save(workout *domain.Workout) error
}
