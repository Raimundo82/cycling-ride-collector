package contracts

import (
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type DailyWorkoutSelector interface {
	GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout
}
