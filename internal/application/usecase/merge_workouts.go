package usecase

import (
	"math"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

func MergeWorkouts(workouts []*domain.Workout, minWorkoutDuration int) *domain.Workout {
	longDurationWorkouts := []*domain.Workout{}
	for _, workout := range workouts {
		if workout.DurationInMin >= minWorkoutDuration {
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
	sumNormalizedPower := 0
	for _, w := range longDurationWorkouts {
		merged.DistanceInKm += w.DistanceInKm
		merged.DurationInMin += w.DurationInMin
		merged.ElevationInMeters += w.ElevationInMeters
		sumAvgPower += w.AvgPowerInWatts * w.DurationInMin
		sumAvgHeartRate += w.AvgHeartRateInBpm * w.DurationInMin
		sumAvgCadence += w.AvgCadenceInRpm * w.DurationInMin
		sumNormalizedPower += w.NormalizedPowerInWatts * w.DurationInMin
		if merged.MaxHeartRateInBpm < w.MaxHeartRateInBpm {
			merged.MaxHeartRateInBpm = w.MaxHeartRateInBpm
		}

	}
	merged.ID = longDurationWorkouts[0].ID
	merged.StartTime = longDurationWorkouts[0].StartTime
	merged.WorkoutType = longDurationWorkouts[0].WorkoutType
	merged.AvgPowerInWatts = weightedAvg(float64(sumAvgPower), merged.DurationInMin)
	merged.AvgHeartRateInBpm = weightedAvg(float64(sumAvgHeartRate), merged.DurationInMin)
	merged.AvgCadenceInRpm = weightedAvg(float64(sumAvgCadence), merged.DurationInMin)
	merged.NormalizedPowerInWatts = weightedAvg(float64(sumNormalizedPower), merged.DurationInMin)

	return merged
}

func weightedAvg(sum float64, totalDuration int) int {
	if totalDuration == 0 {
		return 0
	}
	return int(math.Round(sum / float64(totalDuration)))
}
