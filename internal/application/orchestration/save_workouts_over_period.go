package orchestration

import (
	"errors"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
)

type SaveWorkoutsOrchestrator struct {
	SaveWorkoutUseCase contracts.SaveWorkoutUseCase
}

func (o *SaveWorkoutsOrchestrator) SaveWorkoutsOverPeriod(period Period, minWorkoutDuration int) error {
	var errs []error
	for date := period.StartDate; !date.After(period.EndDate); date = date.AddDate(0, 0, 1) {
		err := o.SaveWorkoutUseCase.Execute(date, minWorkoutDuration)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
