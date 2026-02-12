package usecase

import (
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/samber/lo"
)

var _ contracts.DailyWorkoutPolicy = (*longestWorkout)(nil)

type longestWorkout struct{}

func NewLongestWorkout() contracts.DailyWorkoutPolicy {
	return &longestWorkout{}
}

// GetDailyWorkout implements [contracts.DailyWorkoutPolicy].
func (l *longestWorkout) GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	filteredWorkouts := lo.Filter(dailyWorkouts, func(w *domain.Workout, _ int) bool { return w.DurationInMin >= minWorkoutDuration })
	return lo.MaxBy(filteredWorkouts, func(a, b *domain.Workout) bool { return a.DurationInMin > b.DurationInMin })
}
