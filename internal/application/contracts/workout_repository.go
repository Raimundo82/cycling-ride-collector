package contracts

import "github.com/raimundo82/go-strava-weekly/internal/domain"

type WorkoutSaver interface {
	Save(workout *domain.Workout) error
}

type WorkoutPeriodSaver interface {
	SaveAll(workouts []*domain.Workout) error
}

type WorkoutRepository interface {
	Save(workout *domain.Workout) error
	SaveAll(workouts []*domain.Workout) error
}
