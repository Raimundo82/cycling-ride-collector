package usecase

import (
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type mockDailyWorkoutPolicy struct {
	GetDailyWorkoutCalled int
	Workout               *domain.Workout
}

// GetDailyWorkout implements [contracts.DailyWorkoutPolicy].
func (m *mockDailyWorkoutPolicy) GetDailyWorkout(workouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	m.GetDailyWorkoutCalled++
	return m.Workout
}

var _ contracts.DailyWorkoutPolicy = (*mockDailyWorkoutPolicy)(nil)
