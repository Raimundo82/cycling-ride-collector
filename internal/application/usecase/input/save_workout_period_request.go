package input

import (
	"errors"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type SaveWorkoutPeriodRequest struct {
	Period                 domain.Period
	DailyWorkoutPolicy     string
	MinimalWorkoutDuration int
}

func NewSaveWorkoutPeriodRequest(period domain.Period, dailyWorkoutPolicy string, minimalWorkoutDuration int) (*SaveWorkoutPeriodRequest, error) {
	if period == nil {
		return nil, errors.New("period must be provided")
	}

	if minimalWorkoutDuration < 0 {
		return nil, errors.New("minimal workout duration must be non-negative")
	}

	if dailyWorkoutPolicy != "longest" && dailyWorkoutPolicy != "merge" {
		return nil, errors.New("invalid daily workout policy: allowed values are 'longest' or 'merge'")
	}

	return &SaveWorkoutPeriodRequest{
		Period:                 period,
		DailyWorkoutPolicy:     dailyWorkoutPolicy,
		MinimalWorkoutDuration: minimalWorkoutDuration,
	}, nil
}
