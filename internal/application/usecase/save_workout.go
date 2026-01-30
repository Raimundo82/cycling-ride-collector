package usecase

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

var _ contracts.SaveWorkoutUseCase = (*SaveWorkout)(nil)

type SaveWorkout struct {
	WorkoutRepo     contracts.WorkoutRepository
	WorkoutProvider contracts.WorkoutProvider
}

func (useCase *SaveWorkout) Execute(date time.Time, minWorkoutDuration int) error {
	workouts, err := useCase.WorkoutProvider.GetWorkoutsByDate(date)
	if err != nil {
		return err
	}
	workout := MergeWorkouts(workouts, minWorkoutDuration)

	if workout == nil {
		workout = &domain.Workout{
			StartTime: date,
		}
	}

	return useCase.WorkoutRepo.Save(workout)
}
