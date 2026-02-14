package strava

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

type stubApiClient struct {
	acts             []*ActivityDto
	activitiesErr    error
	wattsStream      *WattsStreamDto
	wattsStreamErr   error
	calls            []int64
	detailedAct      *DetailedActivityDto
	detailedActErr   error
	detailedActCalls []int64
}

var _ StravaApiClient = (*stubApiClient)(nil)

// GetActivitiesByDate implements [StravaApiClient].
func (s *stubApiClient) GetActivitiesByDate(ctx context.Context, d time.Time) ([]*ActivityDto, error) {
	return s.acts, s.activitiesErr
}

// GetActivitiesByPeriod implements [StravaApiClient].
func (s *stubApiClient) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*ActivityDto, error) {
	return s.acts, s.activitiesErr
}

// GetWattsStream implements [StravaApiClient].
func (s *stubApiClient) GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error) {
	s.calls = append(s.calls, activityID)
	return s.wattsStream, s.wattsStreamErr
}

// GetDetailedActivityByID implements [StravaApiClient].
func (s *stubApiClient) GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error) {
	s.detailedActCalls = append(s.detailedActCalls, activityID)
	return s.detailedAct, s.detailedActErr
}

func TestProvider_FiltersAndMapsRides(t *testing.T) {
	Convey("Given a Strava provider", t, func() {
		Convey("When the client returns rides and non-rides", func() {
			stub := &stubApiClient{
				acts: []*ActivityDto{
					{ID: 1, SportType: "Ride", Commute: false},
					{ID: 2, SportType: "Ride", Commute: true},
					{ID: 3, SportType: "MountainBike", Commute: false},
					{ID: 4, SportType: "Run"},
				},
				wattsStream: &WattsStreamDto{WattsData: []int{100, 150, 200, 250, 300, 350, 400, 450}},
				detailedAct: &DetailedActivityDto{ID: 1, LegSensations: "Boas"},
			}

			p := &Provider{apiClient: stub}
			period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
			workouts, err := p.GetWorkoutsByPeriod(period)

			Convey("It should filter only rides", func() {
				So(err, ShouldBeNil)
				So(len(workouts), ShouldEqual, 1)
				So(workouts[0].ID, ShouldEqual, 1)
				So(workouts[0].LegSensations(), ShouldEqual, domain.Good)
				So(stub.calls, ShouldResemble, []int64{1})
				So(stub.detailedActCalls, ShouldResemble, []int64{1})
			})
		})
		Convey("When the client errors", func() {
			stub := &stubApiClient{activitiesErr: errors.New("boom")}
			p := &Provider{apiClient: stub}

			period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
			_, err := p.GetWorkoutsByPeriod(period)

			Convey("It should propagate the error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestProvider_HandlesWattsStreamErrorsGracefully(t *testing.T) {
	Convey("Given a Strava provider where GetWattsStream fails", t, func() {
		stub := &stubApiClient{
			acts: []*ActivityDto{
				{ID: 1, SportType: "Ride", Commute: false, DeviceWatts: true},
				{ID: 2, SportType: "MountainBike", Commute: false},
			},
			wattsStreamErr: errors.New("watts stream unavailable"),
			detailedAct:    &DetailedActivityDto{ID: 1, LegSensations: "Más"},
		}

		p := &Provider{apiClient: stub}
		period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
		workouts, err := p.GetWorkoutsByPeriod(period)

		Convey("It should handle the error gracefully and continue processing", func() {
			So(err, ShouldBeNil)
			So(len(workouts), ShouldEqual, 1)
			So(workouts[0].ID, ShouldEqual, 1)
			So(workouts[0].LegSensations(), ShouldEqual, domain.Bad)
			So(stub.calls, ShouldResemble, []int64{1})
			So(stub.detailedActCalls, ShouldResemble, []int64{1})
		})

		Convey("And the workouts should have zero normalized power", func() {
			So(workouts[0].NormalizedPowerInWatts, ShouldEqual, 0)
		})
	})
}
