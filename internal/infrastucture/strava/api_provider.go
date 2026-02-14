package strava

import (
	"context"
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/samber/lo"
)

var (
	_ contracts.WorkoutProvider = (*ApiProvider)(nil)
)

type ApiProvider struct {
	apiClient StravaApiClient
}

func NewApiProvider(cfg *config.Config, accessToken string) *ApiProvider {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	apiClient := NewStravaApiHttpClient(httpClient, cfg, accessToken)
	return &ApiProvider{apiClient: apiClient}
}

// GetWorkoutsByPeriod implements [contracts.WorkoutProvider].
func (p *ApiProvider) GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error) {
	acts, err := p.apiClient.GetActivitiesByPeriod(context.Background(), period)
	if err != nil {
		return nil, err
	}
	rideActivities := lo.FilterMap(acts, func(a *ActivityDto, _ int) (w *domain.Workout, ok bool) {
		if a.SportType == "Ride" && !a.Commute {
			stream, err := p.apiClient.GetWattsStream(context.Background(), a.ID)
			a.Watts = lo.Ternary(err == nil, stream, &WattsStreamDto{WattsData: []int{}})
			if detailedActivity, err := p.apiClient.GetDetailedActivityByID(context.Background(), a.ID); err == nil && detailedActivity != nil {
				a.LegSensations = detailedActivity.LegSensations
			}
			return MapToWorkout(a), true
		}
		return nil, false
	})
	return rideActivities, nil
}
