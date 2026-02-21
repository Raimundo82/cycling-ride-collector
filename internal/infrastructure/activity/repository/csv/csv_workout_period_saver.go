package activity_csv

import (
	"bytes"
	"encoding/csv"
	"os"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type csvWorkoutPeriodSaver struct {
	filePath string
	buf      *bytes.Buffer
	writer   *csv.Writer
	mapper   func(*domain.Workout) []string
}

// SaveAll implements [contracts.WorkoutRepository].
func (c *csvWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout) error {
	c.buf.Reset()

	for _, workout := range workouts {
		if err := c.writer.Write(c.mapper(workout)); err != nil {
			return err
		}
	}

	c.writer.Flush()
	if err := c.writer.Error(); err != nil {
		return err
	}

	return os.WriteFile(c.filePath, c.buf.Bytes(), 0o644)
}

func NewCSVWorkoutPeriodSaver(filePath string) contracts.WorkoutRepository {
	buf := &bytes.Buffer{}
	return &csvWorkoutPeriodSaver{
		filePath: filePath,
		buf:      buf,
		writer:   csv.NewWriter(buf),
		mapper:   workoutToRecord,
	}
}

var _ contracts.WorkoutRepository = (*csvWorkoutPeriodSaver)(nil)
