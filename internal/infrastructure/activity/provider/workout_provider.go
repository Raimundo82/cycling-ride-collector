package provider

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider/strava"
	"github.com/samber/lo"
)

var _ contracts.WorkoutProvider = (*workoutProvider)(nil)

type workoutProvider struct {
	activityProvider strava.ActivityProvider
}

func NewWorkoutProvider(activityProvider strava.ActivityProvider) *workoutProvider {
	return &workoutProvider{activityProvider: activityProvider}
}

// GetWorkoutsByPeriod implements [contracts.WorkoutProvider].
func (provider *workoutProvider) GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error) {
	activities, err := provider.activityProvider.GetActivitiesByPeriod(context.Background(), period)
	if err != nil {
		return nil, err
	}

	rideActivities := lo.FilterMap(activities, func(a *model.ActivityDto, _ int) (w *domain.Workout, ok bool) {
		if a.SportType == "Ride" && !a.Commute {
			stream, err := provider.activityProvider.GetWattsStream(context.Background(), a.ID)
			a.Watts = lo.Ternary(err == nil, stream, &model.WattsStreamDto{WattsData: []int{}})
			if detailedActivity, err := provider.activityProvider.GetDetailedActivityByID(context.Background(), a.ID); err == nil && detailedActivity != nil {
				a.LegSensations = detailedActivity.LegSensations
			}
			return MapToWorkout(a), true
		}
		return nil, false
	})
	return rideActivities, nil
}
