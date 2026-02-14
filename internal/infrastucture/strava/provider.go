package strava

import (
	"context"
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/samber/lo"
)

var (
	_ contracts.WorkoutProvider = (*Provider)(nil)
	_ contracts.TokenProvider   = (*Provider)(nil)
)

type Provider struct {
	apiClient   StravaApiClient
	oauthClient StravaOAuthClient
	cfg         *config.Config
}

func NewProvider(cfg *config.Config) *Provider {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	apiClient := NewStravaApiHttpClient(httpClient, cfg)
	oauthClient := NewStravaOauthHttpClient(httpClient, cfg)
	return &Provider{apiClient: apiClient, oauthClient: oauthClient, cfg: cfg}
}

// GetWorkoutsByPeriod implements [contracts.WorkoutProvider].
func (p *Provider) GetWorkoutsByPeriod(period domain.Period) ([]*domain.Workout, error) {
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

// RefreshAccessToken implements [contracts.TokenProvider].
func (p *Provider) RefreshAccessToken() (*dto.Token, error) {
	resp, err := p.oauthClient.RefreshAccessToken(&RefreshAccessTokenRequest{})
	if err != nil {
		return nil, err
	}
	return &dto.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    time.Unix(int64(resp.ExpiresAt), 0),
	}, nil
}
