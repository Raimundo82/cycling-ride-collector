package csv

import (
	"os"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type csvWorkoutPeriodSaver struct {
	filePath string
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

func (c *csvWorkoutPeriodSaver) SaveToWriterAll(workouts []*domain.Workout, file *os.File) error {
	panic("unimplemented")
}

func NewCSVWorkoutPeriodSaver(filePath string) contracts.WorkoutRepository {
	return &csvWorkoutPeriodSaver{filePath: filePath}
}

var _ contracts.WorkoutRepository = (*csvWorkoutPeriodSaver)(nil)
