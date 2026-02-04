package usecase

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

var _ contracts.SaveWorkoutUseCase = (*SaveWorkout)(nil)

type SaveWorkout struct {
	dailyWorkout    contracts.DailyWorkoutPolicy
	workoutRepo     contracts.WorkoutRepository
	workoutProvider contracts.WorkoutProvider
}

func NewSaveWorkout(
	dailyWorkout contracts.DailyWorkoutPolicy,
	workoutRepo contracts.WorkoutRepository,
	workoutProvider contracts.WorkoutProvider,
) *SaveWorkout {
	return &SaveWorkout{
		dailyWorkout:    dailyWorkout,
		workoutRepo:     workoutRepo,
		workoutProvider: workoutProvider,
	}
}

func (saveWorkoutUseCase *SaveWorkout) Execute(date time.Time, minWorkoutDuration int) error {
	workouts, err := saveWorkoutUseCase.workoutProvider.GetWorkoutsByDate(date)
	if err != nil {
		return err
	}

	var workout *domain.Workout

	if len(workouts) > 0 {
		workout = saveWorkoutUseCase.dailyWorkout.GetDailyWorkout(workouts, minWorkoutDuration)
	}

	if workout == nil {
		workout = saveWorkoutUseCase.newRestWorkout(date)
	}

	return saveWorkoutUseCase.workoutRepo.Save(workout)
}

func (*SaveWorkout) newRestWorkout(date time.Time) *domain.Workout {
	return domain.NewWorkout(&domain.WorkoutParams{
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
		LegSensations:          "",
	})
}
