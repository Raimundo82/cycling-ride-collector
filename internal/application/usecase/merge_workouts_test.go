package usecase

import (
	"testing"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeWorkouts_EmptyWorkouts(t *testing.T) {
	Convey("Given an empty slice of workouts", t, func() {
		workouts := []*domain.Workout{}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result is an empty workout", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_NoDuration(t *testing.T) {
	Convey("Given a slice with a single workout with no duration", t, func() {
		workout := domain.NewWorkout(
			domain.WorkoutParams{
				Duration: 0,
			})
		workouts := []*domain.Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result should be nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_SingleShortDuration(t *testing.T) {
	Convey("Given a slice with a single workout of short duration", t, func() {
		workout := domain.NewWorkout(
			domain.WorkoutParams{
				Duration: 29,
			})
		workouts := []*domain.Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result should be nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_SingleLongDuration(t *testing.T) {
	Convey("Given a slice with a single workout of long duration", t, func() {
		workout := domain.NewWorkout(
			domain.WorkoutParams{
				Duration: 61,
			})
		workouts := []*domain.Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result is the same workout", func() {
				So(merged, ShouldEqual, workout)
			})
		})
	})
}

func TestMergeWorkouts_MixedShortAndLong(t *testing.T) {
	Convey("Given a slice with mixed short and long duration workouts", t, func() {
		shortWorkout := domain.NewWorkout(
			domain.WorkoutParams{
				Duration: 29,
			})
		longWorkout := domain.NewWorkout(
			domain.WorkoutParams{
				Duration: 75,
			})
		workouts := []*domain.Workout{shortWorkout, longWorkout}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result should be the long workout", func() {
				So(merged, ShouldEqual, longWorkout)
			})
		})
	})
}

func TestMergeWorkouts_AllLongDurations(t *testing.T) {
	Convey("Given a slice with all long duration workouts", t, func() {
		workout1 := domain.NewWorkout(
			domain.WorkoutParams{
				Id:           1,
				WorkoutType:  domain.Estrada,
				StartTime:    "09:00",
				Duration:     50,
				Distance:     20.0,
				Elevation:    100,
				AvgPower:     250,
				AvgHeartRate: 150,
				AvgCadence:   80,
				MaxHeartRate: 150,
			})
		workout2 := domain.NewWorkout(
			domain.WorkoutParams{
				Id:           2,
				WorkoutType:  domain.Estrada,
				StartTime:    "15:00",
				Duration:     100,
				Distance:     30.0,
				Elevation:    200,
				AvgPower:     200,
				AvgHeartRate: 100,
				AvgCadence:   90,
				MaxHeartRate: 180,
			})
		workouts := []*domain.Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.Id, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, "09:00")
				So(merged.WorkoutType, ShouldEqual, domain.Estrada)
				So(merged.Duration, ShouldEqual, 150)
				So(merged.Distance, ShouldEqual, 50.0)
				So(merged.Elevation, ShouldEqual, 300)
				So(merged.AvgPower, ShouldEqual, 217)
				So(merged.AvgHeartRate, ShouldEqual, 117)
				So(merged.AvgCadence, ShouldEqual, 87)
				So(merged.MaxHeartRate, ShouldEqual, 180)
			})
		})
	})
}
