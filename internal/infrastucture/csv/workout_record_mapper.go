package csv

import (
	"fmt"
	"strconv"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

func workoutToRecord(workout *domain.Workout) []string {
	startTime := ""
	if workout.ID >= 0 {
		startTime = workout.StartTime.Format("15:04")
	}

	return []string{
		workout.StartTime.Format("1/2/2006"),
		workout.WorkoutType.String(),
		startTime,
		durationInHoursAndMinutes(workout.DurationInMin),
		floatValueOrEmpty(workout.DistanceInKm),
		intValueOrEmpty(workout.ElevationInMeters),
		intValueOrEmpty(workout.AvgPowerInWatts),
		intValueOrEmpty(workout.NormalizedPowerInWatts),
		intValueOrEmpty(workout.AvgHeartRateInBpm),
		intValueOrEmpty(workout.MaxHeartRateInBpm),
		intValueOrEmpty(workout.AvgCadenceInRpm),
		string(workout.LegSensations()),
	}
}

func intValueOrEmpty(value int) string {
	if value < 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func floatValueOrEmpty(value float64) string {
	if value < 0 {
		return ""
	}
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func durationInHoursAndMinutes(totalMinutes int) string {
	if totalMinutes < 0 {
		return ""
	}
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}
