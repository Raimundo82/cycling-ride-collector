package contracts

import (
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type DailyWorkoutPolicy interface {
	GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout
}
