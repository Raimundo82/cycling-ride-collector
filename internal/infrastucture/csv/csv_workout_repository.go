package csv

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type CSVWorkoutRepository struct {
	filePath string
}

func NewCSVWorkoutRepository(filePath string) *CSVWorkoutRepository {
	return &CSVWorkoutRepository{filePath: filePath}
}

func (r *CSVWorkoutRepository) Save(workout *domain.Workout) (err error) {
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

func (r *CSVWorkoutRepository) SaveToWriter(workout *domain.Workout, w io.Writer) (err error) {
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

func (r *CSVWorkoutRepository) workoutToRecord(workout *domain.Workout) []string {
	startTime := ""
	if workout.ID >= 0 {
		startTime = workout.StartTime.Format("15:04")
	}

	return []string{
		workout.StartTime.Format("1/2/2006"),
		workout.WorkoutType.String(),
		startTime,
		intValueOrEmpty(workout.DurationInMin),
		floatValueOrEmpty(workout.DistanceInKm),
		intValueOrEmpty(workout.ElevationInMeters),
		intValueOrEmpty(workout.AvgPowerInWatts),
		intValueOrEmpty(workout.NormalizedPowerInWatts),
		intValueOrEmpty(workout.AvgHeartRateInBpm),
		intValueOrEmpty(workout.MaxHeartRateInBpm),
		intValueOrEmpty(workout.AvgCadenceInRpm),
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
