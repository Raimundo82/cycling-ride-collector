package usecase

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/samber/lo"
)

type SaveWorkoutPeriodUseCase interface {
	Execute(period domain.Period, minimalWorkoutDuration int) error
}

type saveWorkoutPeriod struct {
	dailyWorkout    contracts.DailyWorkoutPolicy
	workoutRepo     contracts.WorkoutPeriodSaver
	workoutProvider contracts.PeriodWorkoutProvider
}

func NewSaveWorkoutPeriod(
	dailyWorkout contracts.DailyWorkoutPolicy,
	workoutRepo contracts.WorkoutPeriodSaver,
	workoutProvider contracts.PeriodWorkoutProvider,
) SaveWorkoutPeriodUseCase {
	return &saveWorkoutPeriod{
		dailyWorkout:    dailyWorkout,
		workoutRepo:     workoutRepo,
		workoutProvider: workoutProvider,
	}
}

func (saveWorkoutPeriod *saveWorkoutPeriod) Execute(period domain.Period, minimalWorkoutDuration int) error {
	periodWorkouts, err := saveWorkoutPeriod.workoutProvider.GetWorkoutsByPeriod(period)
	if err != nil {
		return err
	}

	workouts := make([]*domain.Workout, 0)
	for date := period.StartDate(); !date.After(period.EndDate()); date = date.AddDate(0, 0, 1) {
		dailyWorkouts := lo.Filter(periodWorkouts, func(w *domain.Workout, _ int) bool {
			return time.Date(w.StartTime.Year(), w.StartTime.Month(), w.StartTime.Day(), 0, 0, 0, 0, w.StartTime.Location()).Equal(date)
		})
		var workout *domain.Workout

		if len(dailyWorkouts) > 0 {
			workout = saveWorkoutPeriod.dailyWorkout.GetDailyWorkout(dailyWorkouts, minimalWorkoutDuration)
		}

		if workout == nil {
			workout = saveWorkoutPeriod.newRestWorkout(date)
		}

		workouts = append(workouts, workout)
	}

	return saveWorkoutPeriod.workoutRepo.SaveAll(workouts)
}

func (*saveWorkoutPeriod) newRestWorkout(date time.Time) *domain.Workout {
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
