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
	return []string{
		workout.StartTime.Format("1/2/2006"),
		workout.WorkoutType.String(),
		valueOrEmpty(workout.StartTime.Format("15:04"), "00:00"),
		valueOrEmpty(strconv.Itoa(workout.DurationInMin), "0"),
		valueOrEmpty(strconv.FormatFloat(workout.DistanceInKm, 'f', 2, 64), "0.00"),
		valueOrEmpty(strconv.Itoa(workout.ElevationInMeters), "0"),
		valueOrEmpty(strconv.Itoa(workout.AvgPowerInWatts), "0"),
		valueOrEmpty(strconv.Itoa(workout.NormalizedPowerInWatts), "0"),
		valueOrEmpty(strconv.Itoa(workout.AvgHeartRateInBpm), "0"),
		valueOrEmpty(strconv.Itoa(workout.MaxHeartRateInBpm), "0"),
		valueOrEmpty(strconv.Itoa(workout.AvgCadenceInRpm), "0"),
	}
}

func valueOrEmpty(value string, zero string) string {
	if value == zero {
		return ""
	}
	return value
}
