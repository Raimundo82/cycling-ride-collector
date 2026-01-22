package usecase

import (
	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
)

type SaveWorkout struct {
	WorkoutRepo     contracts.WorkoutRepository
	WorkoutProvider contracts.WorkoutProvider
}

func (s *SaveWorkout) Execute(unixDate int64) error {
	workout, err := s.WorkoutProvider.FetchWorkout(unixDate)
	if err != nil {
		return err
	}

	if err := s.WorkoutRepo.Save(workout); err != nil {
		return err
	}

	return nil
}
