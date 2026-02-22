package athlete_provider

import (
	"context"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	athlete_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/interfaces"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
)

type athleteDataProvider struct {
	AthleteStatsProvider athlete_interfaces.AthleteStatsProvider
}

// GetAthleteData implements [contracts.AthleteDataProvider].
func (a *athleteDataProvider) GetAthleteData() (*domain.Athlete, error) {
	detailedAthlete, err := a.AthleteStatsProvider.GetDetailedAthlete(context.Background())
	if err != nil {
		return nil, err
	}

	athleteZones, err := a.AthleteStatsProvider.GetAthleteZones(context.Background())
	if err != nil {
		return nil, err
	}

	detailedAthlete.Zones = *athleteZones
	return a.mapAthlete(detailedAthlete), nil
}

var _ contracts.AthleteDataProvider = (*athleteDataProvider)(nil)

func NewAthleteProvider(httpAthleteStatsProvider athlete_interfaces.AthleteStatsProvider) *athleteDataProvider {
	return &athleteDataProvider{
		AthleteStatsProvider: httpAthleteStatsProvider,
	}
}

func (a *athleteDataProvider) mapAthlete(detailedAthlete *model.DetailedAthlete) *domain.Athlete {
	weightInKg := detailedAthlete.Weight
	hrThreshold := detailedAthlete.Zones.HeartRateRangeZones.Zones[1].Max
	pwrThreshold := detailedAthlete.Zones.PowerRangeZones.Zones[1].Max

	return domain.NewAthlete(weightInKg, hrThreshold, pwrThreshold)
}
