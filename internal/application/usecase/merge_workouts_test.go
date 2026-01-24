package usecase

import (
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeWorkouts_EmptyWorkouts(t *testing.T) {
	Convey("Given an empty slice of workouts", t, func() {
		workouts := []*domain.Workout{}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result is nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_NoDuration(t *testing.T) {
	Convey("Given a slice with a single workout with no duration", t, func() {
		workout := domain.NewWorkout(
			domain.WorkoutParams{
				DurationInMin: 0,
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
				DurationInMin: 29,
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
				DurationInMin: 61,
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
				DurationInMin: 29,
			})
		longWorkout := domain.NewWorkout(
			domain.WorkoutParams{
				DurationInMin: 75,
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
				ID:                1,
				WorkoutType:       domain.Estrada,
				StartTime:         time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:     50,
				DistanceInKm:      20.0,
				ElevationInMeters: 100,
				AvgPowerInWatts:   250,
				AvgHeartRateInBpm: 150,
				AvgCadenceInRpm:   80,
				MaxHeartRateInBpm: 150,
			})
		workout2 := domain.NewWorkout(
			domain.WorkoutParams{
				ID:                2,
				WorkoutType:       domain.Estrada,
				StartTime:         time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:     100,
				DistanceInKm:      30.0,
				ElevationInMeters: 200,
				AvgPowerInWatts:   200,
				AvgHeartRateInBpm: 100,
				AvgCadenceInRpm:   90,
				MaxHeartRateInBpm: 180,
			})
		workouts := []*domain.Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := MergeWorkouts(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, domain.Estrada)
				So(merged.DurationInMin, ShouldEqual, 150)
				So(merged.DistanceInKm, ShouldEqual, 50.0)
				So(merged.ElevationInMeters, ShouldEqual, 300)
				So(merged.AvgPowerInWatts, ShouldEqual, 217)
				So(merged.AvgHeartRateInBpm, ShouldEqual, 117)
				So(merged.AvgCadenceInRpm, ShouldEqual, 87)
				So(merged.MaxHeartRateInBpm, ShouldEqual, 180)
			})
		})
	})
}
