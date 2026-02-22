package athlete_strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	athlete_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/interfaces"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
)

const stravaError = "strava error: %s"

type httpAthleteStatsProvider struct {
	httpClient http.Client
	baseUrl    string
}

func NewHttpAthleteStatsProvider(httpClient *http.Client, baseUrl string) *httpAthleteStatsProvider {
	return &httpAthleteStatsProvider{
		baseUrl:    baseUrl,
		httpClient: *httpClient,
	}
}

// GetAthleteZones implements [AthleteStatsProvider].
func (h *httpAthleteStatsProvider) GetAthleteZones(ctx context.Context) (*model.Zones, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseUrl+"/athlete/zones", nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(stravaError, resp.Status)
	}

	var athleteZones model.Zones
	if err := json.NewDecoder(resp.Body).Decode(&athleteZones); err != nil {
		return nil, err
	}

	return &athleteZones, nil
}

// GetDetailedAthlete implements [AthleteStatsProvider].
func (h *httpAthleteStatsProvider) GetDetailedAthlete(ctx context.Context) (*model.DetailedAthlete, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.baseUrl+"/athlete", nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(stravaError, resp.Status)
	}

	var athleteData model.DetailedAthlete
	if err := json.NewDecoder(resp.Body).Decode(&athleteData); err != nil {
		return nil, err
	}

	return &athleteData, nil
}

var _ athlete_interfaces.AthleteStatsProvider = (*httpAthleteStatsProvider)(nil)
