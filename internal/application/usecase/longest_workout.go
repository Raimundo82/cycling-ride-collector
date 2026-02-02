package usecase

import (
	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/samber/lo"
)

var _ contracts.DailyWorkoutPolicy = (*LongestWorkout)(nil)

type LongestWorkout struct{}

// GetDailyWorkout implements [contracts.DailyWorkoutPolicy].
func (l *LongestWorkout) GetDailyWorkout(dailyWorkouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	filteredWorkouts := lo.Filter(dailyWorkouts, func(w *domain.Workout, _ int) bool { return w.DurationInMin >= minWorkoutDuration })
	return lo.MaxBy(filteredWorkouts, func(a, b *domain.Workout) bool { return a.DurationInMin > b.DurationInMin })
}
