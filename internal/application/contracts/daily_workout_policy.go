package contracts

import (
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type DailyWorkoutPolicy interface {
	GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout
}
