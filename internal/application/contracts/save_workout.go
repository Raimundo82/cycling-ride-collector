package contracts

import (
	"time"
)

type SaveWorkoutUseCase interface {
	Execute(date time.Time, minWorkoutDuration int) error
}
