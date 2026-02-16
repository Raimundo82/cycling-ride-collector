package csv

import (
	"fmt"
	"strconv"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/samber/lo"
)

type WorkoutCsvRecordMapper struct{}

func (m *WorkoutCsvRecordMapper) Map(workout *domain.Workout) []string {
	return m.workoutToRecord(workout)
}

func (m *WorkoutCsvRecordMapper) workoutToRecord(workout *domain.Workout) []string {
	return []string{
		workout.StartTime.Format("1/2/2006"),
		workout.WorkoutType.String(),
		lo.Ternary(workout.IsRestDay(), "", workout.StartTime.Format("15:04")),
		formatDurationOrEmpty(workout, workout.DurationInMin),
		formatFloatOrEmpty(workout, workout.DistanceInKm),
		formatIntOrEmpty(workout, workout.ElevationInMeters),
		formatIntOrEmpty(workout, workout.AvgPowerInWatts),
		formatIntOrEmpty(workout, workout.NormalizedPowerInWatts),
		formatIntOrEmpty(workout, workout.AvgHeartRateInBpm),
		formatIntOrEmpty(workout, workout.MaxHeartRateInBpm),
		formatIntOrEmpty(workout, workout.AvgCadenceInRpm),
		string(workout.LegSensations()),
	}
}

func formatIntOrEmpty(workout *domain.Workout, value int) string {
	if workout.IsRestDay() || value < 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func formatFloatOrEmpty(workout *domain.Workout, value float64) string {
	if workout.IsRestDay() || value < 0 {
		return ""
	}
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func formatDurationOrEmpty(workout *domain.Workout, totalMinutes int) string {
	if workout.IsRestDay() || totalMinutes < 0 {
		return ""
	}
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}
