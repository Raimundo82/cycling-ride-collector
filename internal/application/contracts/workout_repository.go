package contracts

import (
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type WorkoutRepository interface {
	SaveAll(workouts []*domain.Workout, athlete *domain.Athlete) error
}
