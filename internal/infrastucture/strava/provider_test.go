package strava

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type stubClient struct {
	acts           []*ActivityDto
	wattsStream    *WattsStreamDto
	activitiesErr  error
	wattsStreamErr error
	calls          []int64
}

var _ client = (*stubClient)(nil)

func (s *stubClient) GetActivitiesByDate(ctx context.Context, d time.Time) ([]*ActivityDto, error) {
	return s.acts, s.activitiesErr
}

func (s *stubClient) GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error) {
	s.calls = append(s.calls, activityID)
	return s.wattsStream, s.wattsStreamErr
}

func TestProvider_FiltersAndMapsRides(t *testing.T) {
	Convey("Given a Strava provider", t, func() {
		Convey("When the client returns rides and non-rides", func() {
			stub := &stubClient{
				acts: []*ActivityDto{
					{ID: 1, Type: "Ride", SportType: "Ride"},
					{ID: 2, Type: "Ride", SportType: "MountainBike"},
					{ID: 3, Type: "Run", SportType: "Run"},
				},
				wattsStream: &WattsStreamDto{WattsData: []int{100, 150, 200, 250, 300, 350, 400, 450}},
			}

			p := NewProvider(stub)
			ws, err := p.GetWorkoutsByDate(time.Now())

			Convey("It should filter only rides", func() {
				So(err, ShouldBeNil)
				So(len(ws), ShouldEqual, 1)
				So(ws[0].ID, ShouldEqual, 1)
				So(stub.calls, ShouldResemble, []int64{1})
			})
		})
		Convey("When the client errors", func() {
			stub := &stubClient{activitiesErr: errors.New("boom")}
			p := NewProvider(stub)

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
			acts: []*ActivityDto{
				{ID: 1, Type: "Ride", SportType: "Ride"},
				{ID: 2, Type: "Ride", SportType: "MountainBike"},
			},
			wattsStreamErr: errors.New("watts stream unavailable"),
		}

		p := NewProvider(stub)
		ws, err := p.GetWorkoutsByDate(time.Now())

		Convey("It should handle the error gracefully and continue processing", func() {
			So(err, ShouldBeNil)
			So(len(ws), ShouldEqual, 1)
			So(ws[0].ID, ShouldEqual, 1)
			So(stub.calls, ShouldResemble, []int64{1})
		})

		Convey("And the workouts should have zero normalized power", func() {
			So(ws[0].NormalizedPowerInWatts, ShouldEqual, 0)
		})
	})
}
