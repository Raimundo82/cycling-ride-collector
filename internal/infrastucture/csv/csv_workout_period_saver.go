package csv

import (
	"encoding/csv"
	"os"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type csvWorkoutPeriodSaver struct {
	filePath string
	mapper   WorkoutCsvRecordMapper
}

// SaveAll implements [contracts.WorkoutPeriodSaver].
func (c *csvWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout) error {
	file, err := os.OpenFile(c.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	return c.SaveToWriterAll(workouts, file)
}

func (c *csvWorkoutPeriodSaver) SaveToWriterAll(workouts []*domain.Workout, file *os.File) (err error) {
	writer := csv.NewWriter(file)
	defer func() {
		writer.Flush()
		if flushErr := writer.Error(); flushErr != nil && err == nil {
			err = flushErr
		}
	}()

	records := make([][]string, 0, len(workouts))
	for _, workout := range workouts {
		record := c.mapper.Map(workout)
		records = append(records, record)
	}
	return writer.WriteAll(records)
}

func NewCSVWorkoutPeriodSaver(filePath string) contracts.WorkoutRepository {
	return &csvWorkoutPeriodSaver{filePath: filePath}
}

var _ contracts.WorkoutRepository = (*csvWorkoutPeriodSaver)(nil)
