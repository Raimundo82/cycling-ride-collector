package usecase

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type SaveWorkoutUseCase interface {
	Execute(date time.Time, minWorkoutDuration int) error
}

var _ SaveWorkoutUseCase = (*saveWorkout)(nil)

type saveWorkout struct {
	dailyWorkout    contracts.DailyWorkoutPolicy
	workoutRepo     contracts.WorkoutSaver
	workoutProvider contracts.SingleWorkoutProvider
}

func NewSaveWorkout(
	dailyWorkout contracts.DailyWorkoutPolicy,
	workoutRepo contracts.WorkoutSaver,
	workoutProvider contracts.SingleWorkoutProvider,
) SaveWorkoutUseCase {
	return &saveWorkout{
		dailyWorkout:    dailyWorkout,
		workoutRepo:     workoutRepo,
		workoutProvider: workoutProvider,
	}
}

func (saveWorkoutUseCase *saveWorkout) Execute(date time.Time, minWorkoutDuration int) error {
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

func (*saveWorkout) newRestWorkout(date time.Time) *domain.Workout {
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
