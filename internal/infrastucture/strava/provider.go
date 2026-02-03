package strava

import (
	"context"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/samber/lo"
)

var _ contracts.WorkoutProvider = (*Provider)(nil)

type Provider struct {
	client Client
}

func NewProvider(c Client) *Provider {
	return &Provider{client: c}
}

func (p *Provider) GetWorkoutsByDate(date time.Time) ([]*domain.Workout, error) {
	acts, err := p.client.GetActivitiesByDate(context.Background(), date)
	if err != nil {
		return nil, err
	}

	rideActivities := lo.FilterMap(acts, func(a *ActivityDto, _ int) (w *domain.Workout, ok bool) {
		if a.SportType == "Ride" && !a.Commute {
			stream, err := p.client.GetWattsStream(context.Background(), a.ID)
			a.Watts = lo.Ternary(err == nil, stream, &WattsStreamDto{WattsData: []int{}})
			detailedActivity, err := p.client.GetDetailedActivityByID(context.Background(), a.ID)
			a.LegSensations = lo.Ternary(err == nil, detailedActivity.LegSensations, "")
			return MapToWorkout(a), true
		}
		return nil, false
	})

	return rideActivities, nil
}
