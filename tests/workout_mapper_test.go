package test

import (
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/strava"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMapToWorkout_With29WattsDataPoints(t *testing.T) {
	Convey("Given a road ride ActivityDto with 29 watts data points", t, func() {
		activity := &strava.ActivityDto{
			ID:                 123456789,
			Type:               "Ride",
			IsTrainer:          false,
			StartDate:          "2024-01-25T10:30:00Z",
			Distance:           25500.4,
			Duration:           3600,
			TotalElevationGain: 500.0,
			AveragePower:       200.0,
			WeightedAvgPower:   220.0,
			AverageHeartRate:   150.0,
			MaxHeartRate:       180.0,
			AverageCadence:     90.0,
			Watts: &strava.WattsStreamDto{WattsData: []int{
				200, 200, 200, 200, 200, 200, 200, 200, 200, 200,
				200, 200, 200, 200, 200, 200, 200, 200, 200, 200,
				200, 200, 200, 200, 200, 200, 200, 200, 200,
			}},
		}

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)

			Convey("Then it should map to Estrada workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Estrada)
			})

			Convey("And it should correctly map all fields", func() {
				So(workout.ID, ShouldEqual, int64(123456789))
				So(workout.WorkoutType, ShouldEqual, domain.Estrada)
				So(workout.StartTime, ShouldEqual, time.Date(2024, 1, 25, 10, 30, 0, 0, time.UTC))
				So(workout.DurationInMin, ShouldEqual, 60)
				So(workout.DistanceInKm, ShouldEqual, 25.5)
				So(workout.ElevationInMeters, ShouldEqual, 500)
				So(workout.AvgPowerInWatts, ShouldEqual, 200)
				So(workout.AvgHeartRateInBpm, ShouldEqual, 150)
				So(workout.MaxHeartRateInBpm, ShouldEqual, 180)
				So(workout.AvgCadenceInRpm, ShouldEqual, 90)
			})
			Convey("And Normalized Power should be zero due to insufficient data points", func() {
				So(workout.NormalizedPowerInWatts, ShouldEqual, 0)
			})
		})
	})
}

func TestMapToWorkout_WithMoreThan30WattsDataPoints(t *testing.T) {
	Convey("Given a trainer ride ActivityDto with more than 30 watts data points", t, func() {
		activity := &strava.ActivityDto{
			ID:                 987654321,
			Type:               "Ride",
			IsTrainer:          true,
			StartDate:          "2024-01-25T18:00:00Z",
			Distance:           30000.0,
			Duration:           5400,
			TotalElevationGain: 250.0,
			AveragePower:       180.0,
			WeightedAvgPower:   200.0,
			AverageHeartRate:   145.0,
			MaxHeartRate:       175.0,
			AverageCadence:     85.0,
			Watts: &strava.WattsStreamDto{WattsData: []int{
				200, 200, 200, 200, 200, 200, 200, 200, 200, 200,
				200, 200, 200, 200, 200, 200, 200, 200, 200, 200,
				200, 200, 200, 200, 200, 200, 200, 200, 200, 200,
			}},
		}

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)

			Convey("Then it should map to Rolo workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Rolo)
			})

			Convey("And it should correctly map all fields", func() {
				So(workout.ID, ShouldEqual, int64(987654321))
				So(workout.WorkoutType, ShouldEqual, domain.Rolo)
				So(workout.StartTime, ShouldEqual, time.Date(2024, 1, 25, 18, 0, 0, 0, time.UTC))
				So(workout.DurationInMin, ShouldEqual, 90)
				So(workout.DistanceInKm, ShouldEqual, 30.0)
				So(workout.ElevationInMeters, ShouldEqual, 250)
				So(workout.AvgPowerInWatts, ShouldEqual, 180)
				So(workout.AvgHeartRateInBpm, ShouldEqual, 145)
				So(workout.MaxHeartRateInBpm, ShouldEqual, 175)
				So(workout.AvgCadenceInRpm, ShouldEqual, 85)
				So(workout.NormalizedPowerInWatts, ShouldEqual, 200)
			})
		})
	})
	Convey("Given an ActivityDto with invalid date format", t, func() {
		activity := &strava.ActivityDto{
			ID:                 111111111,
			Type:               "Ride",
			IsTrainer:          false,
			StartDate:          "invalid-date",
			Distance:           10000.0,
			Duration:           1800,
			TotalElevationGain: 100.0,
			AveragePower:       150.0,
			WeightedAvgPower:   160.0,
			AverageHeartRate:   140.0,
			MaxHeartRate:       170.0,
			AverageCadence:     80.0,
			Watts: &strava.WattsStreamDto{WattsData: []int{
				150, 150, 150, 150, 150, 150, 150, 150, 150, 150,
				150, 150, 150, 150, 150, 150, 150, 150, 150, 150,
				150, 150, 150, 150, 150, 150, 150, 150, 150, 150,
				300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
				300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
				300, 300, 300, 300, 300, 300, 300, 300, 300, 300,
			}},
		}

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)

			Convey("Then StartTime should be zero time", func() {
				So(workout.StartTime, ShouldResemble, time.Time{})
			})

			Convey("And other fields should still be mapped correctly", func() {
				So(workout.ID, ShouldEqual, int64(111111111))
				So(workout.WorkoutType, ShouldEqual, domain.Estrada)
				So(workout.DurationInMin, ShouldEqual, 30)
				So(workout.DistanceInKm, ShouldEqual, 10.0)
				So(workout.NormalizedPowerInWatts, ShouldEqual, 237)
			})
		})
	})
}
