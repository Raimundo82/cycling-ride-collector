package athlete_provider

import (
	"context"
	"fmt"
	"testing"

	athlete_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/interfaces"
	athlete_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/athlete/model"
	. "github.com/smartystreets/goconvey/convey"
)

type httpAthleteStatsProviderMock struct {
	athlete                  *athlete_model.DetailedAthlete
	getDetailedAthleteCalled int
	zones                    *athlete_model.Zones
	getAthleteZonesCalled    int
	detailedAthleteErr       error
	athleteZonesErr          error
}

func (h *httpAthleteStatsProviderMock) GetDetailedAthlete(_ context.Context) (*athlete_model.DetailedAthlete, error) {
	h.getDetailedAthleteCalled++
	if h.detailedAthleteErr != nil {
		return nil, h.detailedAthleteErr
	}
	return h.athlete, nil
}

func (h *httpAthleteStatsProviderMock) GetAthleteZones(_ context.Context) (*athlete_model.Zones, error) {
	h.getAthleteZonesCalled++
	if h.athleteZonesErr != nil {
		return nil, h.athleteZonesErr
	}
	return h.zones, nil
}

var _ athlete_interfaces.AthleteStatsProvider = (*httpAthleteStatsProviderMock)(nil)

func TestGetAthleteDataReturnsAthlete(t *testing.T) {
	Convey("Given a mock HTTP athlete stats provider", t, func() {
		mockAthleteStatsProvider := &httpAthleteStatsProviderMock{
			athlete: &athlete_model.DetailedAthlete{
				Weight: 70,
			},
			zones: &athlete_model.Zones{
				HeartRateRangeZones: &athlete_model.HeartRateRangeZones{Zones: []athlete_model.Zone{{Min: 0, Max: 120}, {Min: 120, Max: 140}, {Min: 140, Max: 155}}},
				PowerRangeZones:     &athlete_model.PowerRangeZones{Zones: []athlete_model.Zone{{Min: 0, Max: 180}, {Min: 180, Max: 240}, {Min: 240, Max: 280}}},
			},
		}

		provider := NewAthleteProvider(mockAthleteStatsProvider)

		Convey("When calling GetAthleteData", func() {
			athlete, err := provider.GetAthleteData()

			Convey("Then it should return an athlete", func() {
				So(err, ShouldBeNil)
				So(athlete, ShouldNotBeNil)
				So(mockAthleteStatsProvider.getDetailedAthleteCalled, ShouldEqual, 1)
				So(mockAthleteStatsProvider.getAthleteZonesCalled, ShouldEqual, 1)
				So(athlete.WeightInKg(), ShouldEqual, 70)
				So(athlete.HeartRateThresholdInBpm(), ShouldEqual, 140)
				So(athlete.PowerThresholdInWatts(), ShouldEqual, 240)
			})
		})
	})
}

func TestGetAthleteDataReturnsErrorWhenAthleteZonesAreIncomplete(t *testing.T) {
	Convey("Given a mock HTTP athlete stats provider", t, func() {
		mockAthleteStatsProvider := &httpAthleteStatsProviderMock{
			athlete: &athlete_model.DetailedAthlete{
				Weight: 70,
			},
			zones: &athlete_model.Zones{
				HeartRateRangeZones: &athlete_model.HeartRateRangeZones{Zones: []athlete_model.Zone{{Min: 0, Max: 120}}},
				PowerRangeZones:     &athlete_model.PowerRangeZones{Zones: []athlete_model.Zone{{Min: 0, Max: 180}}},
			},
		}

		provider := NewAthleteProvider(mockAthleteStatsProvider)

		Convey("When calling GetAthleteData", func() {
			athlete, err := provider.GetAthleteData()

			Convey("Then it should return an athlete", func() {
				So(err, ShouldBeNil)
				So(athlete, ShouldNotBeNil)
				So(mockAthleteStatsProvider.getDetailedAthleteCalled, ShouldEqual, 1)
				So(mockAthleteStatsProvider.getAthleteZonesCalled, ShouldEqual, 1)
				So(athlete.WeightInKg(), ShouldEqual, 70)
				So(athlete.HeartRateThresholdInBpm(), ShouldEqual, 0)
				So(athlete.PowerThresholdInWatts(), ShouldEqual, 0)
			})
		})
	})
}

func TestGetAthleteDataReturnsAthleteDetailedErrorWhenProviderFails(t *testing.T) {
	Convey("Given a mock HTTP athlete stats provider", t, func() {
		mockAthleteStatsProvider := &httpAthleteStatsProviderMock{
			detailedAthleteErr: fmt.Errorf("error getting athlete data"),
		}

		provider := NewAthleteProvider(mockAthleteStatsProvider)

		Convey("When calling GetAthleteData", func() {
			athlete, err := provider.GetAthleteData()

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(athlete, ShouldBeNil)
				So(mockAthleteStatsProvider.getDetailedAthleteCalled, ShouldEqual, 1)
				So(mockAthleteStatsProvider.getAthleteZonesCalled, ShouldEqual, 0)
			})
		})
	})
}

func TestGetAthleteDataReturnsAthleteZonesErrorWhenProviderFails(t *testing.T) {
	Convey("Given a mock HTTP athlete stats provider", t, func() {
		mockAthleteStatsProvider := &httpAthleteStatsProviderMock{
			athleteZonesErr: fmt.Errorf("error getting athlete zones"),
		}

		provider := NewAthleteProvider(mockAthleteStatsProvider)

		Convey("When calling GetAthleteData", func() {
			athlete, err := provider.GetAthleteData()

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(athlete, ShouldBeNil)
				So(mockAthleteStatsProvider.getDetailedAthleteCalled, ShouldEqual, 1)
				So(mockAthleteStatsProvider.getAthleteZonesCalled, ShouldEqual, 1)
			})
		})
	})
}
