package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type csvWorkoutSaver struct {
	filePath string
}

var (
	_ contracts.WorkoutSaver       = (*csvWorkoutSaver)(nil)
	_ contracts.WorkoutPeriodSaver = (*csvWorkoutSaver)(nil)
)

func NewCSVWorkoutSaver(filePath string) contracts.WorkoutSaver {
	return &csvWorkoutSaver{filePath: filePath}
}

func (r *csvWorkoutSaver) Save(workout *domain.Workout) (err error) {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return r.SaveToWriter(workout, file)
}

func (r *csvWorkoutSaver) SaveAll(workouts []*domain.Workout) (err error) {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return r.SaveToWriterAll(workouts, file)
}

func (r *csvWorkoutSaver) SaveToWriter(workout *domain.Workout, w io.Writer) (err error) {
	writer := csv.NewWriter(w)
	defer func() {
		writer.Flush()
		if flushErr := writer.Error(); flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	record := r.workoutToRecord(workout)
	return writer.Write(record)
}

func (r *csvWorkoutSaver) SaveToWriterAll(workouts []*domain.Workout, w io.Writer) (err error) {
	writer := csv.NewWriter(w)
	defer func() {
		writer.Flush()
		if flushErr := writer.Error(); flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	records := make([][]string, 0, len(workouts))
	for _, workout := range workouts {
		record := r.workoutToRecord(workout)
		records = append(records, record)
	}
	return writer.WriteAll(records)
}

func (r *csvWorkoutSaver) workoutToRecord(workout *domain.Workout) []string {
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
