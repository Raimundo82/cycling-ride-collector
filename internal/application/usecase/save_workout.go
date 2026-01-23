package usecase

import (
	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
)

type SaveWorkout struct {
	WorkoutRepo     contracts.WorkoutRepository
	WorkoutProvider contracts.WorkoutProvider
}

func (useCase *SaveWorkout) Execute(unixDate int64, minWorkoutDuration int) error {
	workouts, err := useCase.WorkoutProvider.GetWorkoutsByDate(unixDate)
	if err != nil {
		return err
	}

	workout := MergeWorkouts(workouts, minWorkoutDuration)

	if workout == nil {
		return nil
	}

	return useCase.WorkoutRepo.Save(workout)
}
