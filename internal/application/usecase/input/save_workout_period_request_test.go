package input

import (
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSaveWorkoutPeriodRequest_GivenValidPeriodAndMinimalWorkoutDurationAndPolicies_WhenNewSaveWorkoutPeriodRequestIsInvoked_ThenValidRequestIsReturned(t *testing.T) {
	Convey("Given a valid period, minimal workout duration and daily workout policy", t, func() {
		p, _ := domain.NewPeriod(
			time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		)
		minDuration := 30
		policy := "longest"

		Convey("When NewSaveWorkoutPeriodRequest is invoked", func() {
			r, err := NewSaveWorkoutPeriodRequest(p, policy, minDuration)

			Convey("Then a valid SaveWorkoutPeriodRequest is returned", func() {
				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
			})
		})
	})
}

func TestNewSaveWorkoutPeriodRequest_GivenInvalidPeriodAndMinimalWorkoutDurationAndPolicies_WhenNewSaveWorkoutPeriodRequestIsInvoked_ThenErrorIsReturned(t *testing.T) {
	Convey("Given an invalid period, minimal workout duration and daily workout policy", t, func() {
		var p domain.Period
		minDuration := 30
		policy := "longest"

		Convey("When NewSaveWorkoutPeriodRequest is invoked", func() {
			r, err := NewSaveWorkoutPeriodRequest(p, policy, minDuration)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(r, ShouldBeNil)
			})
		})
	})
}

func TestNewSaveWorkoutPeriodRequest_GivenInvalidMinimalWorkoutDuration_WhenNewSaveWorkoutPeriodRequestIsInvoked_ThenErrorIsReturned(t *testing.T) {
	Convey("Given an invalid period, minimal workout duration and daily workout policy", t, func() {
		p, _ := domain.NewPeriod(
			time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		)
		minDuration := -1
		policy := "longest"

		Convey("When NewSaveWorkoutPeriodRequest is invoked", func() {
			r, err := NewSaveWorkoutPeriodRequest(p, policy, minDuration)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(r, ShouldBeNil)
			})
		})
	})
}

func TestNewSaveWorkoutPeriodRequest_GivenInvalidPolicy_WhenNewSaveWorkoutPeriodRequestIsInvoked_ThenErrorIsReturned(t *testing.T) {
	Convey("Given an invalid period, minimal workout duration and daily workout policy", t, func() {
		p, _ := domain.NewPeriod(
			time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		)
		minDuration := 30
		policy := "invalid_policy"

		Convey("When NewSaveWorkoutPeriodRequest is invoked", func() {
			r, err := NewSaveWorkoutPeriodRequest(p, policy, minDuration)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(r, ShouldBeNil)
			})
		})
	})
}
