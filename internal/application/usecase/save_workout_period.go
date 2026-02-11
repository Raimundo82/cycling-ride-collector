package usecase

import (
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type SaveWorkoutPeriodUseCase interface {
	Execute(period domain.Period, minimalWorkoutDuration int) error
}

type saveWorkoutPeriod struct {
	dailyWorkout    contracts.DailyWorkoutPolicy
	workoutRepo     contracts.WorkoutRepository
	workoutProvider contracts.PeriodWorkoutProvider
}

func NewSaveWorkoutPeriod(
	dailyWorkout contracts.DailyWorkoutPolicy,
	workoutRepo contracts.WorkoutRepository,
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

	workoutsByDate := saveWorkoutPeriod.groupWorkoutsByDate(periodWorkouts)
	workouts := make([]*domain.Workout, 0)
	for date := period.StartDate(); !date.After(period.EndDate()); date = date.AddDate(0, 0, 1) {
		normalizedDate := normalizeToDate(date)
		dailyWorkouts := workoutsByDate[normalizedDate]

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

func (saveWorkoutPeriod *saveWorkoutPeriod) groupWorkoutsByDate(workouts []*domain.Workout) map[time.Time][]*domain.Workout {
	workoutsByDate := make(map[time.Time][]*domain.Workout)
	for _, workout := range workouts {
		date := normalizeToDate(workout.StartTime)
		workoutsByDate[date] = append(workoutsByDate[date], workout)
	}
	return workoutsByDate
}

func normalizeToDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
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
