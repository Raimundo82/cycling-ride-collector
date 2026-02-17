package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider/strava"
	. "github.com/smartystreets/goconvey/convey"
)

type mockActivityProvider struct {
	acts             []*model.ActivityDto
	activitiesErr    error
	wattsStream      *model.WattsStreamDto
	wattsStreamErr   error
	calls            []int64
	detailedAct      *model.DetailedActivityDto
	detailedActErr   error
	detailedActCalls []int64
}

var _ strava.ActivityProvider = (*mockActivityProvider)(nil)

// GetActivitiesByPeriod implements [strava.ActivityProvider].
func (s *mockActivityProvider) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*model.ActivityDto, error) {
	return s.acts, s.activitiesErr
}

// GetWattsStream implements [strava.ActivityProvider].
func (s *mockActivityProvider) GetWattsStream(ctx context.Context, activityID int64) (*model.WattsStreamDto, error) {
	s.calls = append(s.calls, activityID)
	return s.wattsStream, s.wattsStreamErr
}

// GetDetailedActivityByID implements [strava.ActivityProvider].
func (s *mockActivityProvider) GetDetailedActivityByID(ctx context.Context, activityID int64) (*model.DetailedActivityDto, error) {
	s.detailedActCalls = append(s.detailedActCalls, activityID)
	return s.detailedAct, s.detailedActErr
}

func TestWorkoutProviderReturnsFilteredAndMappedWorkouts(t *testing.T) {
	Convey("Given a activity provider", t, func() {
		activityProvider := &mockActivityProvider{
			acts: []*model.ActivityDto{
				{ID: 1, SportType: "Ride", Commute: false},
				{ID: 2, SportType: "Ride", Commute: true},
				{ID: 3, SportType: "MountainBike", Commute: false},
				{ID: 4, SportType: "Run"},
			},
			wattsStream: &model.WattsStreamDto{WattsData: []int{100, 150, 200, 250, 300, 350, 400, 450}},
			detailedAct: &model.DetailedActivityDto{ID: 1, LegSensations: "Boas"},
		}

		workoutProvider := &workoutProvider{activityProvider: activityProvider}

		Convey("When getting workouts by period", func() {
			period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
			workouts, err := workoutProvider.GetWorkoutsByPeriod(period)

			Convey("It should filter and map only rides", func() {
				So(err, ShouldBeNil)
				So(len(workouts), ShouldEqual, 1)
				So(workouts[0].ID, ShouldEqual, 1)
				So(workouts[0].LegSensations(), ShouldEqual, domain.Good)
				So(activityProvider.calls, ShouldResemble, []int64{1})
				So(activityProvider.detailedActCalls, ShouldResemble, []int64{1})
			})
		})
	})
}

func TestWorkoutProviderReturnsErrorWhenReturnedByActivityProvider(t *testing.T) {
	Convey("Given a activity provider", t, func() {
		activityProvider := &mockActivityProvider{activitiesErr: errors.New("boom")}
		workoutProvider := &workoutProvider{activityProvider: activityProvider}

		Convey("When the client has errors", func() {
			period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
			_, err := workoutProvider.GetWorkoutsByPeriod(period)

			Convey("Then it should propagate the error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "boom")
				So(activityProvider.calls, ShouldBeEmpty)
				So(activityProvider.detailedActCalls, ShouldBeEmpty)
			})
		})
	})
}

func TestWorkoutProviderHandlesWattsStreamErrorsGracefully(t *testing.T) {
	Convey("Given a Strava provider where GetWattsStream fails", t, func() {
		activityProvider := &mockActivityProvider{
			acts: []*model.ActivityDto{
				{ID: 1, SportType: "Ride", Commute: false, DeviceWatts: true},
				{ID: 2, SportType: "MountainBike", Commute: false},
			},
			wattsStreamErr: errors.New("watts stream unavailable"),
			detailedAct:    &model.DetailedActivityDto{ID: 1, LegSensations: "Más"},
		}

		workoutProvider := &workoutProvider{activityProvider: activityProvider}
		Convey("When getting workouts by period", func() {
			period, _ := domain.NewPeriod(time.Now(), time.Now().Add(24*time.Hour))
			workouts, err := workoutProvider.GetWorkoutsByPeriod(period)

			Convey("Then It should handle the error gracefully and continue processing", func() {
				So(err, ShouldBeNil)
				So(len(workouts), ShouldEqual, 1)
				So(workouts[0].ID, ShouldEqual, 1)
				So(workouts[0].LegSensations(), ShouldEqual, domain.Bad)
				So(activityProvider.calls, ShouldResemble, []int64{1})
				So(activityProvider.detailedActCalls, ShouldResemble, []int64{1})
			})

			Convey("Then the workouts should have zero normalized power", func() {
				So(workouts[0].NormalizedPowerInWatts, ShouldEqual, 0)
			})
		})
	})
}

func TestNewWorkoutProviderReturnsProviderWithInjectedActivityProvider(t *testing.T) {
	Convey("Given an activity provider", t, func() {
		activityProvider := &mockActivityProvider{}

		Convey("When creating a new workout provider", func() {
			workoutProvider := NewWorkoutProvider(activityProvider)

			Convey("Then it should return a provider with injected dependency", func() {
				So(workoutProvider, ShouldNotBeNil)
				So(workoutProvider.activityProvider, ShouldEqual, activityProvider)
			})
		})
	})
}
