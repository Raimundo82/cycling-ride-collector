package usecase

import "github.com/raimundo82/go-strava-weekly/internal/application/contracts"

type SaveWorkout struct {
	WorkoutRepo     contracts.WorkoutRepository
	WorkoutProvider contracts.WorkoutProvider
}

func (s *SaveWorkout) Execute() error {
	workouts, err := s.WorkoutProvider.FetchWorkouts()
	if err != nil {
		return err
	}

	for _, workout := range workouts {
		if err := s.WorkoutRepo.Save(workout); err != nil {
			return err
		}
	}

	return nil
}
