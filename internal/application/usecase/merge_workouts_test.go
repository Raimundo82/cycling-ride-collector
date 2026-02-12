package usecase

import (
	"testing"
	"time"

	. "github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeWorkouts_EmptyWorkouts(t *testing.T) {
	Convey("Given an empty slice of workouts", t, func() {
		workouts := []*Workout{}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_NoDuration(t *testing.T) {
	Convey("Given a slice with a single workout with no duration", t, func() {
		workout := NewWorkout(&WorkoutParams{
			DurationInMin: 0,
		})
		workouts := []*Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result should be nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_SingleShortDuration(t *testing.T) {
	Convey("Given a slice with a single workout of short duration", t, func() {
		workout := NewWorkout(
			&WorkoutParams{
				DurationInMin: 29,
			})
		workouts := []*Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result should be nil", func() {
				So(merged, ShouldBeNil)
			})
		})
	})
}

func TestMergeWorkouts_SingleLongDuration(t *testing.T) {
	Convey("Given a slice with a single workout of long duration", t, func() {
		workout := NewWorkout(
			&WorkoutParams{
				DurationInMin: 61,
			})
		workouts := []*Workout{workout}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is the same workout", func() {
				So(merged, ShouldEqual, workout)
			})
		})
	})
}

func TestMergeWorkouts_MixedShortAndLong(t *testing.T) {
	Convey("Given a slice with mixed short and long duration workouts", t, func() {
		shortWorkout := NewWorkout(
			&WorkoutParams{
				DurationInMin: 29,
			})
		longWorkout := NewWorkout(
			&WorkoutParams{
				DurationInMin: 75,
			})
		workouts := []*Workout{shortWorkout, longWorkout}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result should be the long workout", func() {
				So(merged, ShouldEqual, longWorkout)
			})
		})
	})
}

func TestMergeWorkouts_AllLongDurations(t *testing.T) {
	Convey("Given a slice with all long duration workouts", t, func() {
		workout1 := NewWorkout(
			&WorkoutParams{
				ID:                     1,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:          50,
				DistanceInKm:           20.0,
				ElevationInMeters:      100,
				AvgPowerInWatts:        250,
				AvgHeartRateInBpm:      150,
				AvgCadenceInRpm:        80,
				MaxHeartRateInBpm:      150,
				NormalizedPowerInWatts: 280,
			})
		workout2 := NewWorkout(
			&WorkoutParams{
				ID:                     2,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        200,
				AvgHeartRateInBpm:      100,
				AvgCadenceInRpm:        90,
				NormalizedPowerInWatts: 225,
				MaxHeartRateInBpm:      180,
			})
		workouts := []*Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, Estrada)
				So(merged.DurationInMin, ShouldEqual, 150)
				So(merged.DistanceInKm, ShouldEqual, 50.0)
				So(merged.ElevationInMeters, ShouldEqual, 300)
				So(merged.AvgPowerInWatts, ShouldEqual, 217)
				So(merged.AvgHeartRateInBpm, ShouldEqual, 117)
				So(merged.AvgCadenceInRpm, ShouldEqual, 87)
				So(merged.MaxHeartRateInBpm, ShouldEqual, 180)
				So(merged.NormalizedPowerInWatts, ShouldEqual, 243)
			})
		})
	})
}

func TestMergeWorkouts_WithMissingHeartRateAndPowerData(t *testing.T) {
	Convey("Given a slice with all long duration workouts with some missing heart rate and power data", t, func() {
		workout1 := NewWorkout(
			&WorkoutParams{
				ID:                     1,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:          50,
				DistanceInKm:           20.0,
				ElevationInMeters:      100,
				AvgPowerInWatts:        -1,
				AvgHeartRateInBpm:      -1,
				AvgCadenceInRpm:        80,
				MaxHeartRateInBpm:      -1,
				NormalizedPowerInWatts: -1,
			})
		workout2 := NewWorkout(
			&WorkoutParams{
				ID:                     2,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        200,
				AvgHeartRateInBpm:      100,
				AvgCadenceInRpm:        90,
				NormalizedPowerInWatts: 225,
				MaxHeartRateInBpm:      180,
			})
		workouts := []*Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, Estrada)
				So(merged.DurationInMin, ShouldEqual, 150)
				So(merged.DistanceInKm, ShouldEqual, 50.0)
				So(merged.ElevationInMeters, ShouldEqual, 300)
				So(merged.AvgPowerInWatts, ShouldEqual, 200)
				So(merged.AvgHeartRateInBpm, ShouldEqual, 100)
				So(merged.AvgCadenceInRpm, ShouldEqual, 87)
				So(merged.MaxHeartRateInBpm, ShouldEqual, 180)
				So(merged.NormalizedPowerInWatts, ShouldEqual, 225)
			})
		})
	})
}

func TestMergeWorkouts_WithMissingHeartRateAndPowerDataInAllWorkouts(t *testing.T) {
	Convey("Given a slice with all long duration workouts with some missing heart rate and power data", t, func() {
		workout1 := NewWorkout(
			&WorkoutParams{
				ID:                     1,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:          50,
				DistanceInKm:           20.0,
				ElevationInMeters:      100,
				AvgPowerInWatts:        -1,
				AvgHeartRateInBpm:      -1,
				AvgCadenceInRpm:        -1,
				MaxHeartRateInBpm:      -1,
				NormalizedPowerInWatts: -1,
			})
		workout2 := NewWorkout(
			&WorkoutParams{
				ID:                     2,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        -1,
				AvgHeartRateInBpm:      -1,
				AvgCadenceInRpm:        -1,
				NormalizedPowerInWatts: -1,
				MaxHeartRateInBpm:      -1,
			})
		workouts := []*Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, Estrada)
				So(merged.DurationInMin, ShouldEqual, 150)
				So(merged.DistanceInKm, ShouldEqual, 50.0)
				So(merged.ElevationInMeters, ShouldEqual, 300)
				So(merged.AvgPowerInWatts, ShouldEqual, -1)
				So(merged.AvgHeartRateInBpm, ShouldEqual, -1)
				So(merged.AvgCadenceInRpm, ShouldEqual, -1)
				So(merged.MaxHeartRateInBpm, ShouldEqual, -1)
				So(merged.NormalizedPowerInWatts, ShouldEqual, -1)
			})
		})
	})
}

func TestMergeShortAndLongWorkouts_WithMissingHeartRateAndPowerDataInAllWorkouts(t *testing.T) {
	Convey("Given a slice with short and long duration workouts with some missing heart rate and power data", t, func() {
		workout1 := NewWorkout(
			&WorkoutParams{
				ID:                     1,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:          20,
				DistanceInKm:           20.0,
				ElevationInMeters:      100,
				AvgPowerInWatts:        -1,
				AvgHeartRateInBpm:      -1,
				AvgCadenceInRpm:        -1,
				MaxHeartRateInBpm:      -1,
				NormalizedPowerInWatts: -1,
			})
		workout2 := NewWorkout(
			&WorkoutParams{
				ID:                     2,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        150,
				AvgHeartRateInBpm:      120,
				AvgCadenceInRpm:        90,
				NormalizedPowerInWatts: 180,
				MaxHeartRateInBpm:      140,
			})
		workouts := []*Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 2)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, Estrada)
				So(merged.DurationInMin, ShouldEqual, 100)
				So(merged.DistanceInKm, ShouldEqual, 30.0)
				So(merged.ElevationInMeters, ShouldEqual, 200)
				So(merged.AvgPowerInWatts, ShouldEqual, 150)
				So(merged.AvgHeartRateInBpm, ShouldEqual, 120)
				So(merged.AvgCadenceInRpm, ShouldEqual, 90)
				So(merged.MaxHeartRateInBpm, ShouldEqual, 140)
				So(merged.NormalizedPowerInWatts, ShouldEqual, 180)
			})
		})
	})
}

func TestMergeLongWorkouts_WithNoMissingHeartRateAndPowerDataInAllWorkouts(t *testing.T) {
	Convey("Given a slice with long duration workouts with no missing heart rate and power data", t, func() {
		workout1 := NewWorkout(
			&WorkoutParams{
				ID:                     1,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        0,
				AvgHeartRateInBpm:      120,
				AvgCadenceInRpm:        90,
				MaxHeartRateInBpm:      140,
				NormalizedPowerInWatts: 0,
			})
		workout2 := NewWorkout(
			&WorkoutParams{
				ID:                     2,
				WorkoutType:            Estrada,
				StartTime:              time.Date(2023, time.January, 1, 15, 0, 0, 0, time.UTC),
				DurationInMin:          100,
				DistanceInKm:           30.0,
				ElevationInMeters:      200,
				AvgPowerInWatts:        150,
				AvgHeartRateInBpm:      120,
				AvgCadenceInRpm:        90,
				MaxHeartRateInBpm:      140,
				NormalizedPowerInWatts: 180,
			})
		workouts := []*Workout{workout1, workout2}

		Convey("When MergeWorkouts is called", func() {
			merged := NewMergeWorkouts().GetDailyWorkout(workouts, 30)

			Convey("Then the result is a merged workout with summed duration and distance", func() {
				So(merged.ID, ShouldEqual, 1)
				So(merged.StartTime, ShouldEqual, time.Date(2023, time.January, 1, 9, 0, 0, 0, time.UTC))
				So(merged.WorkoutType, ShouldEqual, Estrada)
				So(merged.DurationInMin, ShouldEqual, 200)
				So(merged.DistanceInKm, ShouldEqual, 60.0)
				So(merged.ElevationInMeters, ShouldEqual, 400)
				So(merged.AvgPowerInWatts, ShouldEqual, 75)
				So(merged.AvgHeartRateInBpm, ShouldEqual, 120)
				So(merged.AvgCadenceInRpm, ShouldEqual, 90)
				So(merged.MaxHeartRateInBpm, ShouldEqual, 140)
				So(merged.NormalizedPowerInWatts, ShouldEqual, 90)
			})
		})
	})
}
