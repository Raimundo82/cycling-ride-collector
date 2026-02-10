package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/strava"
	. "github.com/smartystreets/goconvey/convey"
)

type stubClient struct {
	acts             []*strava.ActivityDto
	activitiesErr    error
	wattsStream      *strava.WattsStreamDto
	wattsStreamErr   error
	calls            []int64
	detailedAct      *strava.DetailedActivityDto
	detailedActErr   error
	detailedActCalls []int64
}

var _ strava.Client = (*stubClient)(nil)

// GetActivitiesByDate implements [strava.Client].
func (s *stubClient) GetActivitiesByDate(ctx context.Context, d time.Time) ([]*strava.ActivityDto, error) {
	return s.acts, s.activitiesErr
}

// GetActivitiesByPeriod implements [strava.Client].
func (s *stubClient) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*strava.ActivityDto, error) {
	panic("unimplemented")
}

// GetWattsStream implements [strava.Client].
func (s *stubClient) GetWattsStream(ctx context.Context, activityID int64) (*strava.WattsStreamDto, error) {
	s.calls = append(s.calls, activityID)
	return s.wattsStream, s.wattsStreamErr
}

// GetDetailedActivityByID implements [strava.Client].
func (s *stubClient) GetDetailedActivityByID(ctx context.Context, activityID int64) (*strava.DetailedActivityDto, error) {
	s.detailedActCalls = append(s.detailedActCalls, activityID)
	return s.detailedAct, s.detailedActErr
}

func TestProvider_FiltersAndMapsRides(t *testing.T) {
	Convey("Given a Strava provider", t, func() {
		Convey("When the client returns rides and non-rides", func() {
			stub := &stubClient{
				acts: []*strava.ActivityDto{
					{ID: 1, SportType: "Ride", Commute: false},
					{ID: 2, SportType: "Ride", Commute: true},
					{ID: 3, SportType: "MountainBike", Commute: false},
					{ID: 4, SportType: "Run"},
				},
				wattsStream: &strava.WattsStreamDto{WattsData: []int{100, 150, 200, 250, 300, 350, 400, 450}},
				detailedAct: &strava.DetailedActivityDto{ID: 1, LegSensations: "Boas"},
			}

			p := strava.NewProvider(stub)
			workouts, err := p.GetWorkoutsByDate(time.Now())

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
			stub := &stubClient{activitiesErr: errors.New("boom")}
			p := strava.NewProvider(stub)

			_, err := p.GetWorkoutsByDate(time.Now())

			Convey("It should propagate the error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestProvider_HandlesWattsStreamErrorsGracefully(t *testing.T) {
	Convey("Given a Strava provider where GetWattsStream fails", t, func() {
		stub := &stubClient{
			acts: []*strava.ActivityDto{
				{ID: 1, SportType: "Ride", Commute: false, DeviceWatts: true},
				{ID: 2, SportType: "MountainBike", Commute: false},
			},
			wattsStreamErr: errors.New("watts stream unavailable"),
			detailedAct:    &strava.DetailedActivityDto{ID: 1, LegSensations: "Más"},
		}

		p := strava.NewProvider(stub)
		workouts, err := p.GetWorkoutsByDate(time.Now())

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
