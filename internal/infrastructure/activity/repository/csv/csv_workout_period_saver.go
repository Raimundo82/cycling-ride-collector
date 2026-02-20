package activity_csv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type csvWorkoutPeriodSaver struct {
	filePath string
	mapper   WorkoutCsvRecordMapper
}

// SaveAll implements [contracts.WorkoutRepository].
func (c *csvWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout) (err error) {
	file, err := os.OpenFile(c.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	return c.SaveToWriterAll(workouts, file)
}

func (c *csvWorkoutPeriodSaver) SaveToWriterAll(workouts []*domain.Workout, w io.Writer) (err error) {
	writer := csv.NewWriter(w)

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
