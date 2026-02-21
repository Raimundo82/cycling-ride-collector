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
	mapper   *WorkoutCsvRecordMapper
}

// SaveAll implements [contracts.WorkoutRepository].
func (c *csvWorkoutPeriodSaver) SaveAll(workouts []*domain.Workout) (err error) {
	c.buf.Reset()

	for _, workout := range workouts {
		if err = c.writer.Write(c.mapper.Map(workout)); err != nil {
			return err
		}
	}

	c.writer.Flush()
	if err = c.writer.Error(); err != nil {
		return err
	}

	file, err := os.OpenFile(c.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	_, err = file.Write(c.buf.Bytes())
	return err
}

func NewCSVWorkoutPeriodSaver(filePath string) contracts.WorkoutRepository {
	buf := &bytes.Buffer{}
	return &csvWorkoutPeriodSaver{
		filePath: filePath,
		buf:      buf,
		writer:   csv.NewWriter(buf),
		mapper:   &WorkoutCsvRecordMapper{},
	}
}

var _ contracts.WorkoutRepository = (*csvWorkoutPeriodSaver)(nil)
