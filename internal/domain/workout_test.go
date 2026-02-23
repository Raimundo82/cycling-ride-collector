package domain

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewWorkoutReturnsAValidWorkoutWithValidInputParameters(t *testing.T) {
	Convey("Given valid input parameters", t, func() {
		workoutParams := WorkoutParams{
			ID:                     1,
			WorkoutType:            Estrada,
			StartTime:              time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			DistanceInKm:           100.5,
			DurationInMin:          120,
			ElevationInMeters:      1500,
			AvgPowerInWatts:        250,
			NormalizedPowerInWatts: 260,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      170,
			AvgCadenceInRpm:        90,
			LegSensations:          "Boas",
		}

		Convey("When NewWorkout is invoked", func() {
			workout := NewWorkout(workoutParams)

			Convey("Then a valid Workout is returned", func() {
				So(workout, ShouldNotBeNil)
				So(workout.ID, ShouldEqual, workoutParams.ID)
				So(workout.WorkoutType, ShouldEqual, workoutParams.WorkoutType)
				So(workout.StartTime, ShouldEqual, workoutParams.StartTime)
				So(workout.DistanceInKm, ShouldEqual, workoutParams.DistanceInKm)
				So(workout.DurationInMin, ShouldEqual, workoutParams.DurationInMin)
				So(workout.ElevationInMeters, ShouldEqual, workoutParams.ElevationInMeters)
				So(workout.AvgPowerInWatts, ShouldEqual, workoutParams.AvgPowerInWatts)
				So(workout.NormalizedPowerInWatts, ShouldEqual, workoutParams.NormalizedPowerInWatts)
				So(workout.AvgHeartRateInBpm, ShouldEqual, workoutParams.AvgHeartRateInBpm)
				So(workout.MaxHeartRateInBpm, ShouldEqual, workoutParams.MaxHeartRateInBpm)
				So(workout.AvgCadenceInRpm, ShouldEqual, workoutParams.AvgCadenceInRpm)
				So(workout.legSensations, ShouldEqual, Good)
			})
		})
	})
}

func TestNewWorkoutReturnsAValidWorkoutWithNoParams(t *testing.T) {
	Convey("Given a empty workout params", t, func() {
		workoutParams := WorkoutParams{}

		Convey("When NewWorkout is invoked", func() {
			workout := NewWorkout(workoutParams)

			Convey("Then a valid Workout is returned with default leg sensations", func() {
				So(workout, ShouldNotBeNil)
				So(workout.WorkoutType, ShouldEqual, Descanso)
				So(string(workout.legSensations), ShouldEqual, "")
			})
		})
	})
}

func TestWorkoutReturnErrorWhenNilInputIsPassed(t *testing.T) {
	Convey("Given a nil workout params", t, func() {
		Convey("When NewWorkout is invoked", func() {
			var workoutParams WorkoutParams
			workout := NewWorkout(workoutParams)

			Convey("Then a valid Workout is returned with default leg sensations", func() {
				So(workout, ShouldNotBeNil)
				So(workout.WorkoutType, ShouldEqual, Descanso)
				So(string(workout.legSensations), ShouldEqual, "")
			})
		})
	})
}

func TestWorkoutTypeString(t *testing.T) {
	Convey("Given a WorkoutType", t, func() {
		Convey("When String is invoked", func() {
			So(Estrada.String(), ShouldEqual, "Estrada")
			So(Rolo.String(), ShouldEqual, "Rolo")
			So(Prova.String(), ShouldEqual, "Prova")
			So(Descanso.String(), ShouldEqual, "Descanso")
		})
	})
}

func TestWorkoutSetLegSensationReturns(t *testing.T) {
	Convey("Given Workout with legSensations string", t, func() {
		workout := &Workout{}
		Convey("When SetLegSensations is invoked", func() {
			workout.SetLegSensations("Excelentes")
			So(workout.legSensations, ShouldEqual, Excellent)
			workout.SetLegSensations("Muito Boas")
			So(workout.legSensations, ShouldEqual, VeryGood)
			workout.SetLegSensations("Boas")
			So(workout.legSensations, ShouldEqual, Good)
			workout.SetLegSensations("Médias")
			So(workout.legSensations, ShouldEqual, Medium)
			workout.SetLegSensations("Más")
			So(workout.legSensations, ShouldEqual, Bad)
			workout.SetLegSensations("Muito Más")
			So(workout.legSensations, ShouldEqual, VeryBad)
			workout.SetLegSensations("")
			So(string(workout.legSensations), ShouldEqual, "")
		})
	})
}

func TestWorkoutLegSensationsReturnsTheCurrentLegSensations(t *testing.T) {
	Convey("Given Workout with legSensations set", t, func() {
		workout := &Workout{}
		workout.legSensations = Good

		Convey("When LegSensations is invoked", func() {
			So(workout.LegSensations(), ShouldEqual, Good)
		})
	})
}

func TestIsRestDayReturnsTrueWhenWorkoutTypeIsDescanso(t *testing.T) {
	Convey("Given a Workout with WorkoutType Descanso", t, func() {
		workout := &Workout{
			WorkoutType: Descanso,
		}

		Convey("When IsRestDay is invoked", func() {
			So(workout.IsRestDay(), ShouldBeTrue)
		})
	})
}

func TestIsRestDayReturnsFalseWhenWorkoutTypeIsNotDescanso(t *testing.T) {
	Convey("Given a Workout with WorkoutType Estrada", t, func() {
		workout := &Workout{
			WorkoutType: Estrada,
		}

		Convey("When IsRestDay is invoked", func() {
			So(workout.IsRestDay(), ShouldBeFalse)
		})
	})
}
