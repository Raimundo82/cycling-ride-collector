package activity_csv

import (
	"bytes"
	"encoding/csv"
	"os"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/samber/lo"
)

type csvWorkoutPeriodSaver struct {
	filePath string
	buf      *bytes.Buffer
	writer   *csv.Writer
}

// SaveAll implements [contracts.WorkoutRepository].
func (c *csvWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout, athlete *domain.Athlete) error {
	c.buf.Reset()
	records := lo.Map(workouts, func(w *domain.Workout, _ int) []string {
		return workoutToRecord(w, athlete.WeightInKg())
	})

	if err := c.writer.WriteAll(records); err != nil {
		return err
	}

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
	}
}

var _ contracts.WorkoutRepository = (*csvWorkoutPeriodSaver)(nil)
