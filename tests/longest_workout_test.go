package test

import (
	"testing"

	. "github.com/raimundo82/go-strava-weekly/internal/application/usecase"
	. "github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLongestWorkout_ForSingleWorkoutAboveMinDuration(t *testing.T) {
	Convey("Given an long duration workout", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{ID: 1, DurationInMin: 90}),
		}

		Convey("When LongestWorkout is called", func() {
			longest := (&LongestWorkout{}).GetDailyWorkout(workouts, 30)
			Convey("Then the result is the longest workout", func() {
				So(longest, ShouldNotBeNil)
				So(longest.ID, ShouldEqual, 1)
				So(longest.DurationInMin, ShouldEqual, 90)
			})
		})
	})
}

func TestLongestWorkout_ForSingleWorkoutUnderMinDuration(t *testing.T) {
	Convey("Given an short duration workout", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{ID: 1, DurationInMin: 29}),
		}

		Convey("When LongestWorkout is called", func() {
			longest := (&LongestWorkout{}).GetDailyWorkout(workouts, 30)
			Convey("Then the result is nil", func() {
				So(longest, ShouldBeNil)
			})
		})
	})
}

func TestLongestWorkout_ForMultipleWorkoutsShortAndLongDuration(t *testing.T) {
	Convey("Given an short and long duration workout", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{ID: 1, DurationInMin: 29}),
			NewWorkout(&WorkoutParams{ID: 2, DurationInMin: 90}),
		}

		Convey("When LongestWorkout is called", func() {
			longest := (&LongestWorkout{}).GetDailyWorkout(workouts, 30)
			Convey("Then the result is workout with Id 2", func() {
				So(longest, ShouldNotBeNil)
				So(longest.ID, ShouldEqual, 2)
				So(longest.DurationInMin, ShouldEqual, 90)
			})
		})
	})
}

func TestLongestWorkout_ForMultipleWorkoutsOfShortDuration(t *testing.T) {
	Convey("Given short duration workouts", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{ID: 1, DurationInMin: 29}),
			NewWorkout(&WorkoutParams{ID: 2, DurationInMin: 20}),
		}

		Convey("When LongestWorkout is called", func() {
			longest := (&LongestWorkout{}).GetDailyWorkout(workouts, 30)
			Convey("Then the result is nil", func() {
				So(longest, ShouldBeNil)
			})
		})
	})
}

func TestLongestWorkout_ForMultipleWorkoutsOfLongDuration(t *testing.T) {
	Convey("Given long duration workouts", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{ID: 1, DurationInMin: 120}),
			NewWorkout(&WorkoutParams{ID: 2, DurationInMin: 90}),
		}

		Convey("When LongestWorkout is called", func() {
			longest := (&LongestWorkout{}).GetDailyWorkout(workouts, 30)
			Convey("Then the result is the longest workout", func() {
				So(longest, ShouldNotBeNil)
				So(longest.ID, ShouldEqual, 1)
				So(longest.DurationInMin, ShouldEqual, 120)
			})
		})
	})
}
