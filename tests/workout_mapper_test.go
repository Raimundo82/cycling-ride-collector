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
			IsTrainer:          false,
			WorkoutType:        10,
			StartDate:          "2024-01-25T10:30:00Z",
			Distance:           25500.4,
			Duration:           3600,
			TotalElevationGain: 500.0,
			AveragePower:       200.0,
			WeightedAvgPower:   220.0,
			AverageHeartRate:   150.0,
			MaxHeartRate:       180.0,
			AverageCadence:     90.0,
			HasHeartRate:       true,
			DeviceWatts:        true,
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
			IsTrainer:          true,
			WorkoutType:        10,
			StartDate:          "2024-01-25T18:00:00Z",
			Distance:           30000.0,
			Duration:           5400,
			TotalElevationGain: 250.0,
			AveragePower:       180.0,
			WeightedAvgPower:   200.0,
			AverageHeartRate:   145.0,
			MaxHeartRate:       175.0,
			AverageCadence:     85.0,
			HasHeartRate:       true,
			DeviceWatts:        true,
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
			IsTrainer:          false,
			WorkoutType:        12,
			StartDate:          "invalid-date",
			Distance:           10000.0,
			Duration:           1800,
			TotalElevationGain: 100.0,
			AveragePower:       150.0,
			WeightedAvgPower:   160.0,
			AverageHeartRate:   140.0,
			MaxHeartRate:       170.0,
			AverageCadence:     80.0,
			HasHeartRate:       true,
			DeviceWatts:        true,
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

func TestMapToWorkoutDifferentWorkoutTypes(t *testing.T) {
	isTrainer := func(a *strava.ActivityDto) { a.IsTrainer = true }
	isNotTrainer := func(a *strava.ActivityDto) { a.IsTrainer = false }
	NoneWorkout := func(a *strava.ActivityDto) { a.WorkoutType = 10 }
	Workout := func(a *strava.ActivityDto) { a.WorkoutType = 12 }
	ProvaWorkout := func(a *strava.ActivityDto) { a.WorkoutType = 11 }
	HasHeartRate := func(a *strava.ActivityDto) { a.HasHeartRate = true }
	NoHeartRate := func(a *strava.ActivityDto) { a.HasHeartRate = false }
	HasDeviceWatts := func(a *strava.ActivityDto) { a.DeviceWatts = true }
	NoDeviceWatts := func(a *strava.ActivityDto) { a.DeviceWatts = false }

	Convey("Given ActivityDto with workout type 10 (None)", t, func() {
		activity := newActivity(isNotTrainer, NoneWorkout)

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Estrada workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Estrada)
			})
		})
	})
	Convey("Given ActivityDto with workout type 11 (Race)", t, func() {
		activity := newActivity(isNotTrainer, ProvaWorkout)

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Prova workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Prova)
			})
		})
	})

	Convey("Given ActivityDto with workout type 12 (Workout)", t, func() {
		activity := newActivity(isNotTrainer, Workout)

		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Estrada workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Estrada)
			})
		})
	})

	Convey("Given ActivityDto with IsTrainer true and workout type 12 (Workout)", t, func() {
		activity := newActivity(isTrainer, Workout)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Rolo workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Rolo)
			})
		})
	})

	Convey("Given ActivityDto with IsTrainer true and workout type 11 (Prova)", t, func() {
		activity := newActivity(isTrainer, ProvaWorkout)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Rolo workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Rolo)
			})
		})
	})

	Convey("Given ActivityDto with IsTrainer true and workout type 10 (None)", t, func() {
		activity := newActivity(isTrainer, NoneWorkout)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should map to Rolo workout type", func() {
				So(workout.WorkoutType, ShouldEqual, domain.Rolo)
			})
		})
	})

	Convey("Given ActivityDto without hearRate and deviceWatts data", t, func() {
		activity := newActivity(NoHeartRate, NoDeviceWatts)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should have sentinel values for missing data", func() {
				So(workout.AvgHeartRateInBpm, ShouldEqual, -1)
				So(workout.MaxHeartRateInBpm, ShouldEqual, -1)
				So(workout.AvgPowerInWatts, ShouldEqual, -1)
				So(workout.NormalizedPowerInWatts, ShouldEqual, -1)
			})
		})
	})

	Convey("Given ActivityDto with hearRate and deviceWatts data", t, func() {
		activity := newActivity(HasHeartRate, HasDeviceWatts)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should have actual values for present data", func() {
				So(workout.AvgHeartRateInBpm, ShouldEqual, 0)
				So(workout.MaxHeartRateInBpm, ShouldEqual, 0)
				So(workout.AvgPowerInWatts, ShouldEqual, 0)
				So(workout.NormalizedPowerInWatts, ShouldEqual, 0)
			})
		})
	})

	Convey("Given ActivityDto with avg cadence of 0", t, func() {
		activity := newActivity(HasHeartRate, HasDeviceWatts)
		Convey("When MapToWorkout is called", func() {
			workout := strava.MapToWorkout(activity)
			Convey("Then it should have actual values for present data", func() {
				So(workout.AvgCadenceInRpm, ShouldEqual, -1)
			})
		})
	})
}

func newActivity(opts ...func(*strava.ActivityDto)) *strava.ActivityDto {
	activity := &strava.ActivityDto{
		ID:        222222222,
		StartDate: "2024-01-26T09:00:00Z",
		Watts:     &strava.WattsStreamDto{WattsData: []int{}},
	}
	for _, opt := range opts {
		opt(activity)
	}
	return activity
}
