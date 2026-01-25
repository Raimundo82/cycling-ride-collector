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
	client client
}

func NewProvider(c client) *Provider {
	return &Provider{client: c}
}

func (p *Provider) GetWorkoutsByDate(data time.Time) ([]*domain.Workout, error) {
	acts, err := p.client.GetActivityByDate(context.Background(), data)
	if err != nil {
		return nil, err
	}

	rideActivities := lo.FilterMap(acts, func(a *ActivityDto, _ int) (w *domain.Workout, ok bool) {
		if a.Type == "Ride" && a.SportType == "Ride" {
			stream, err := p.client.GetWattsStream(context.Background(), a.ID)
			a.Watts = lo.Ternary(err == nil, stream, nil)
			return MapToWorkout(a), true
		}
		return nil, false
	})

	return rideActivities, nil
}
