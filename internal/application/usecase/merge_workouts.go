package usecase

import (
	"math"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

func MergeWorkouts(workouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	longDurationWorkouts := []*domain.Workout{}
	for _, workout := range workouts {
		if workout.Duration >= minWorkoutDuration {
			longDurationWorkouts = append(longDurationWorkouts, workout)
		}
	}

	if len(longDurationWorkouts) == 0 {
		return nil
	}

	if len(longDurationWorkouts) == 1 {
		return longDurationWorkouts[0]
	}

	merged := &domain.Workout{}
	sumAvgPower := 0
	sumAvgHeartRate := 0
	sumAvgCadence := 0
	for _, w := range longDurationWorkouts {
		merged.Distance += w.Distance
		merged.Duration += w.Duration
		merged.Elevation += w.Elevation
		sumAvgPower += w.AvgPower * w.Duration
		sumAvgHeartRate += w.AvgHeartRate * w.Duration
		sumAvgCadence += w.AvgCadence * w.Duration
		if merged.MaxHeartRate < w.MaxHeartRate {
			merged.MaxHeartRate = w.MaxHeartRate
		}

	}
	merged.Id = longDurationWorkouts[0].Id
	merged.StartTime = longDurationWorkouts[0].StartTime
	merged.WorkoutType = domain.Estrada
	merged.AvgPower = weightedAvg(float64(sumAvgPower), merged.Duration)
	merged.AvgHeartRate = weightedAvg(float64(sumAvgHeartRate), merged.Duration)
	merged.AvgCadence = weightedAvg(float64(sumAvgCadence), merged.Duration)

	return merged
}

func weightedAvg(sum float64, totalDuration int) int {
	if totalDuration == 0 {
		return 0
	}
	return int(math.Round(sum / float64(totalDuration)))
}
