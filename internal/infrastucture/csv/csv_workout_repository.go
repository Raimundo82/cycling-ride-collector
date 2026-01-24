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

func (r *CSVWorkoutRepository) Save(workout *domain.Workout) error {
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

func (r *CSVWorkoutRepository) SaveToWriter(workout *domain.Workout, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	record := r.workoutToRecord(workout)
	return writer.Write(record)
}

func (r *CSVWorkoutRepository) workoutToRecord(workout *domain.Workout) []string {
	return []string{
		workout.WorkoutType.String(),
		workout.StartTime.Format("15:04"),
		strconv.Itoa(workout.DurationInMin),
		strconv.FormatFloat(workout.DistanceInKm, 'f', 2, 64),
		strconv.Itoa(workout.ElevationInMeters),
		strconv.Itoa(workout.AvgPowerInWatts),
		strconv.Itoa(workout.NormalizedPowerInWatts),
		strconv.Itoa(workout.AvgHeartRateInBpm),
		strconv.Itoa(workout.MaxHeartRateInBpm),
		strconv.Itoa(workout.AvgCadenceInRpm),
	}
}
