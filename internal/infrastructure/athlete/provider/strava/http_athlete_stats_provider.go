package athlete_strava

import (
	"net/http"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	athlete_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/interfaces"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
)

type httpAthleteStatsProvider struct {
	Config     *config.Config
	HttpClient http.Client
}

func NewHttpAthleteStatsProvider(httpClient http.Client, cfg *config.Config) *httpAthleteStatsProvider {
	return &httpAthleteStatsProvider{
		Config:     cfg,
		HttpClient: httpClient,
	}
}

// GetAthleteZones implements [AthleteStatsProvider].
func (h *httpAthleteStatsProvider) GetAthleteZones() (*model.Zones, error) {
	panic("unimplemented")
}

// GetDetailedAthlete implements [AthleteStatsProvider].
func (h *httpAthleteStatsProvider) GetDetailedAthlete() (*model.DetailedAthlete, error) {
	panic("unimplemented")
}

var _ athlete_interfaces.AthleteStatsProvider = (*httpAthleteStatsProvider)(nil)
