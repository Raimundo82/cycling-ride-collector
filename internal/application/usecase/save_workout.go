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
			ID:                     -1,
			StartTime:              date,
			WorkoutType:            domain.Descanso,
			DistanceInKm:           -1.00,
			DurationInMin:          -1,
			ElevationInMeters:      -1,
			AvgPowerInWatts:        -1,
			NormalizedPowerInWatts: -1,
			AvgHeartRateInBpm:      -1,
			MaxHeartRateInBpm:      -1,
			AvgCadenceInRpm:        -1,
		}
	}

	return useCase.WorkoutRepo.Save(workout)
}
