package usecase

import (
	"math"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/samber/lo"
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

	for _, w := range longDurationWorkouts {
		merged.DistanceInKm += w.DistanceInKm
		merged.DurationInMin += w.DurationInMin
		merged.ElevationInMeters += w.ElevationInMeters

	}
	merged.ID = longDurationWorkouts[0].ID
	merged.StartTime = longDurationWorkouts[0].StartTime
	merged.WorkoutType = longDurationWorkouts[0].WorkoutType
	merged.AvgPowerInWatts = mergeMetric(longDurationWorkouts, func(w *domain.Workout) int { return w.AvgPowerInWatts })
	merged.AvgHeartRateInBpm = mergeMetric(longDurationWorkouts, func(w *domain.Workout) int { return w.AvgHeartRateInBpm })
	merged.AvgCadenceInRpm = mergeMetric(longDurationWorkouts, func(w *domain.Workout) int { return w.AvgCadenceInRpm })
	merged.NormalizedPowerInWatts = mergeMetric(longDurationWorkouts, func(w *domain.Workout) int { return w.NormalizedPowerInWatts })
	merged.MaxHeartRateInBpm = lo.Max(lo.Map(longDurationWorkouts, func(w *domain.Workout, _ int) int { return w.MaxHeartRateInBpm }))
	return merged
}

func weightedAvg(sum float64, totalDuration int) int {
	if totalDuration == 0 {
		return -1
	}
	return int(math.Round(sum / float64(totalDuration)))
}

func mergeMetric(workouts []*domain.Workout, metric func(*domain.Workout) int) int {
	sum := 0
	duration := 0
	for _, w := range workouts {
		val := metric(w)
		if val != -1 {
			sum += val * w.DurationInMin
			duration += w.DurationInMin
		}
	}
	return weightedAvg(float64(sum), duration)
}
