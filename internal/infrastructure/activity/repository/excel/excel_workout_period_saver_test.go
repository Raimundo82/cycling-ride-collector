package activity_excel

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/xuri/excelize/v2"
)

var testAthlete = NewAthlete(70, 135, 240)

func TestExcelWorkoutPeriodSaverSavesWorkoutsToExcelFile(t *testing.T) {
	Convey("Given an array of workouts, an athlete and an Excel file", t, func() {
		tmpDir, _ := os.MkdirTemp("", "excel-test")
		defer func() { _ = os.RemoveAll(tmpDir) }()
		filePath := filepath.Join(tmpDir, "test_workouts.xlsx")

		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType:            Estrada,
				StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
				DurationInMin:          60,
				DistanceInKm:           25.5,
				ElevationInMeters:      500,
				AvgPowerInWatts:        200,
				NormalizedPowerInWatts: 220,
				AvgHeartRateInBpm:      150,
				MaxHeartRateInBpm:      180,
				AvgCadenceInRpm:        90,
				LegSensations:          string(Good),
			}),
			NewWorkout(WorkoutParams{
				WorkoutType:            Estrada,
				StartTime:              time.Date(2024, 6, 2, 9, 0, 0, 0, time.UTC),
				DurationInMin:          45,
				DistanceInKm:           20.0,
				ElevationInMeters:      300,
				AvgPowerInWatts:        180,
				NormalizedPowerInWatts: 190,
				AvgHeartRateInBpm:      140,
				MaxHeartRateInBpm:      170,
				AvgCadenceInRpm:        85,
				LegSensations:          string(Medium),
			}),
		}

		saver := NewExcelWorkoutPeriodSaver(filePath)

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the Excel file should contain the correct data", func() {
				So(err, ShouldBeNil)
				f, openErr := excelize.OpenFile(filePath)
				So(openErr, ShouldBeNil)
				defer func() { _ = f.Close() }()

				val, _ := f.GetCellValue(defaultSheet, "A1")
				So(val, ShouldEqual, "6/1/2024")
				val, _ = f.GetCellValue(defaultSheet, "B1")
				So(val, ShouldEqual, "Estrada")
				val, _ = f.GetCellValue(defaultSheet, "C1")
				So(val, ShouldEqual, "10:30")
				val, _ = f.GetCellValue(defaultSheet, "D1")
				So(val, ShouldEqual, "1h0m")
				val, _ = f.GetCellValue(defaultSheet, "E1")
				So(val, ShouldEqual, "25.50")
				val, _ = f.GetCellValue(defaultSheet, "M1")
				So(val, ShouldEqual, "70.00")

				val, _ = f.GetCellValue(defaultSheet, "A2")
				So(val, ShouldEqual, "6/2/2024")
				val, _ = f.GetCellValue(defaultSheet, "B2")
				So(val, ShouldEqual, "Estrada")
				val, _ = f.GetCellValue(defaultSheet, "C2")
				So(val, ShouldEqual, "09:00")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverSavesWeeklyWorkoutsMonToSun(t *testing.T) {
	Convey("Given a full week of workouts from Monday to Sunday", t, func() {
		tmpDir, _ := os.MkdirTemp("", "excel-test")
		defer func() { _ = os.RemoveAll(tmpDir) }()
		filePath := filepath.Join(tmpDir, "test_workouts.xlsx")

		// Week of June 3-9, 2024 (Mon-Sun)
		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType: Estrada,
				StartTime:   time.Date(2024, 6, 3, 7, 0, 0, 0, time.UTC),
				DurationInMin: 60, DistanceInKm: 30.0, ElevationInMeters: 400,
				AvgPowerInWatts: 200, NormalizedPowerInWatts: 210,
				AvgHeartRateInBpm: 145, MaxHeartRateInBpm: 175, AvgCadenceInRpm: 88,
				LegSensations: string(Good),
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Rolo,
				StartTime:   time.Date(2024, 6, 4, 18, 0, 0, 0, time.UTC),
				DurationInMin: 45, DistanceInKm: 20.0, ElevationInMeters: 0,
				AvgPowerInWatts: 190, NormalizedPowerInWatts: 200,
				AvgHeartRateInBpm: 140, MaxHeartRateInBpm: 165, AvgCadenceInRpm: 90,
				LegSensations: string(Medium),
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Descanso,
				StartTime:   time.Date(2024, 6, 5, 0, 0, 0, 0, time.UTC),
				DurationInMin: -1, DistanceInKm: -1, ElevationInMeters: -1,
				AvgPowerInWatts: -1, NormalizedPowerInWatts: -1,
				AvgHeartRateInBpm: -1, MaxHeartRateInBpm: -1, AvgCadenceInRpm: -1,
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Estrada,
				StartTime:   time.Date(2024, 6, 6, 6, 30, 0, 0, time.UTC),
				DurationInMin: 90, DistanceInKm: 45.0, ElevationInMeters: 800,
				AvgPowerInWatts: 210, NormalizedPowerInWatts: 225,
				AvgHeartRateInBpm: 150, MaxHeartRateInBpm: 180, AvgCadenceInRpm: 85,
				LegSensations: string(VeryGood),
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Descanso,
				StartTime:   time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC),
				DurationInMin: -1, DistanceInKm: -1, ElevationInMeters: -1,
				AvgPowerInWatts: -1, NormalizedPowerInWatts: -1,
				AvgHeartRateInBpm: -1, MaxHeartRateInBpm: -1, AvgCadenceInRpm: -1,
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Estrada,
				StartTime:   time.Date(2024, 6, 8, 8, 0, 0, 0, time.UTC),
				DurationInMin: 120, DistanceInKm: 60.0, ElevationInMeters: 1000,
				AvgPowerInWatts: 195, NormalizedPowerInWatts: 215,
				AvgHeartRateInBpm: 148, MaxHeartRateInBpm: 178, AvgCadenceInRpm: 87,
				LegSensations: string(Good),
			}),
			NewWorkout(WorkoutParams{
				WorkoutType: Descanso,
				StartTime:   time.Date(2024, 6, 9, 0, 0, 0, 0, time.UTC),
				DurationInMin: -1, DistanceInKm: -1, ElevationInMeters: -1,
				AvgPowerInWatts: -1, NormalizedPowerInWatts: -1,
				AvgHeartRateInBpm: -1, MaxHeartRateInBpm: -1, AvgCadenceInRpm: -1,
			}),
		}

		saver := NewExcelWorkoutPeriodSaver(filePath)

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the Excel file should have 7 rows for Mon-Sun", func() {
				So(err, ShouldBeNil)
				f, openErr := excelize.OpenFile(filePath)
				So(openErr, ShouldBeNil)
				defer func() { _ = f.Close() }()

				// Monday row
				val, _ := f.GetCellValue(defaultSheet, "A1")
				So(val, ShouldEqual, "6/3/2024")
				val, _ = f.GetCellValue(defaultSheet, "B1")
				So(val, ShouldEqual, "Estrada")

				// Tuesday row
				val, _ = f.GetCellValue(defaultSheet, "A2")
				So(val, ShouldEqual, "6/4/2024")
				val, _ = f.GetCellValue(defaultSheet, "B2")
				So(val, ShouldEqual, "Rolo")

				// Wednesday rest day
				val, _ = f.GetCellValue(defaultSheet, "A3")
				So(val, ShouldEqual, "6/5/2024")
				val, _ = f.GetCellValue(defaultSheet, "B3")
				So(val, ShouldEqual, "Descanso")
				val, _ = f.GetCellValue(defaultSheet, "C3")
				So(val, ShouldEqual, "")

				// Thursday row
				val, _ = f.GetCellValue(defaultSheet, "A4")
				So(val, ShouldEqual, "6/6/2024")

				// Friday rest day
				val, _ = f.GetCellValue(defaultSheet, "A5")
				So(val, ShouldEqual, "6/7/2024")
				val, _ = f.GetCellValue(defaultSheet, "B5")
				So(val, ShouldEqual, "Descanso")

				// Saturday row
				val, _ = f.GetCellValue(defaultSheet, "A6")
				So(val, ShouldEqual, "6/8/2024")
				val, _ = f.GetCellValue(defaultSheet, "B6")
				So(val, ShouldEqual, "Estrada")

				// Sunday rest day
				val, _ = f.GetCellValue(defaultSheet, "A7")
				So(val, ShouldEqual, "6/9/2024")
				val, _ = f.GetCellValue(defaultSheet, "B7")
				So(val, ShouldEqual, "Descanso")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverWithCustomStartCell(t *testing.T) {
	Convey("Given workouts and a custom start cell B3", t, func() {
		tmpDir, _ := os.MkdirTemp("", "excel-test")
		defer func() { _ = os.RemoveAll(tmpDir) }()
		filePath := filepath.Join(tmpDir, "test_workouts.xlsx")

		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType:            Estrada,
				StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
				DurationInMin:          60,
				DistanceInKm:           25.5,
				ElevationInMeters:      500,
				AvgPowerInWatts:        200,
				NormalizedPowerInWatts: 220,
				AvgHeartRateInBpm:      150,
				MaxHeartRateInBpm:      180,
				AvgCadenceInRpm:        90,
				LegSensations:          string(Good),
			}),
		}

		saver := NewExcelWorkoutPeriodSaverWithOptions(filePath, defaultSheet, "B3")

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the data should start at cell B3", func() {
				So(err, ShouldBeNil)
				f, openErr := excelize.OpenFile(filePath)
				So(openErr, ShouldBeNil)
				defer func() { _ = f.Close() }()

				val, _ := f.GetCellValue(defaultSheet, "A1")
				So(val, ShouldEqual, "")
				val, _ = f.GetCellValue(defaultSheet, "B3")
				So(val, ShouldEqual, "6/1/2024")
				val, _ = f.GetCellValue(defaultSheet, "C3")
				So(val, ShouldEqual, "Estrada")
				val, _ = f.GetCellValue(defaultSheet, "D3")
				So(val, ShouldEqual, "10:30")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverCreatesFileIfNotExist(t *testing.T) {
	Convey("Given a workout, an athlete and a non-existing Excel file", t, func() {
		filePath := "test_create_workouts.xlsx"
		defer func() { _ = os.Remove(filePath) }()
		saver := NewExcelWorkoutPeriodSaver(filePath)

		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType:            Estrada,
				StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
				DurationInMin:          60,
				DistanceInKm:           25.5,
				ElevationInMeters:      500,
				AvgPowerInWatts:        200,
				NormalizedPowerInWatts: 220,
				AvgHeartRateInBpm:      150,
				MaxHeartRateInBpm:      180,
				AvgCadenceInRpm:        90,
				LegSensations:          string(VeryGood),
			}),
		}

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the file should be created with the correct data", func() {
				So(err, ShouldBeNil)
				f, openErr := excelize.OpenFile(filePath)
				So(openErr, ShouldBeNil)
				defer func() { _ = f.Close() }()

				val, _ := f.GetCellValue(defaultSheet, "A1")
				So(val, ShouldEqual, "6/1/2024")
				val, _ = f.GetCellValue(defaultSheet, "L1")
				So(val, ShouldEqual, "Muito Boas")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverReturnsErrorGivenAnInvalidFilePath(t *testing.T) {
	Convey("Given an invalid file path", t, func() {
		invalidPath := "/invalid/path/that/does/not/exist/workouts.xlsx"
		saver := NewExcelWorkoutPeriodSaver(invalidPath)
		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType:            Estrada,
				StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
				DurationInMin:          60,
				DistanceInKm:           25.5,
				ElevationInMeters:      500,
				AvgPowerInWatts:        200,
				NormalizedPowerInWatts: 220,
				AvgHeartRateInBpm:      150,
				MaxHeartRateInBpm:      180,
				AvgCadenceInRpm:        90,
				LegSensations:          string(VeryGood),
			}),
		}

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverReturnsNoErrorGivenWorkoutWithSentinelValues(t *testing.T) {
	Convey("Given a workout with sentinel values (rest day)", t, func() {
		tmpDir, _ := os.MkdirTemp("", "excel-test")
		defer func() { _ = os.RemoveAll(tmpDir) }()
		filePath := filepath.Join(tmpDir, "test_workouts.xlsx")

		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				ID:                     -1,
				StartTime:              time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				WorkoutType:            Descanso,
				DistanceInKm:           -1.00,
				DurationInMin:          -1,
				ElevationInMeters:      -1,
				AvgPowerInWatts:        -1,
				NormalizedPowerInWatts: -1,
				AvgHeartRateInBpm:      -1,
				MaxHeartRateInBpm:      -1,
				AvgCadenceInRpm:        -1,
				LegSensations:          "",
			}),
		}

		saver := NewExcelWorkoutPeriodSaver(filePath)

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the file should contain the rest day with empty values", func() {
				f, openErr := excelize.OpenFile(filePath)
				So(openErr, ShouldBeNil)
				defer func() { _ = f.Close() }()

				val, _ := f.GetCellValue(defaultSheet, "A1")
				So(val, ShouldEqual, "6/1/2024")
				val, _ = f.GetCellValue(defaultSheet, "B1")
				So(val, ShouldEqual, "Descanso")
				val, _ = f.GetCellValue(defaultSheet, "C1")
				So(val, ShouldEqual, "")
				val, _ = f.GetCellValue(defaultSheet, "D1")
				So(val, ShouldEqual, "")
				val, _ = f.GetCellValue(defaultSheet, "M1")
				So(val, ShouldEqual, "70.00")
			})
		})
	})
}

func TestExcelWorkoutPeriodSaverReturnsErrorGivenAnInvalidStartCell(t *testing.T) {
	Convey("Given an invalid start cell", t, func() {
		tmpDir, _ := os.MkdirTemp("", "excel-test")
		defer func() { _ = os.RemoveAll(tmpDir) }()
		filePath := filepath.Join(tmpDir, "test_workouts.xlsx")

		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				WorkoutType: Estrada,
				StartTime:   time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
				DurationInMin: 60, DistanceInKm: 25.5, ElevationInMeters: 500,
				AvgPowerInWatts: 200, NormalizedPowerInWatts: 220,
				AvgHeartRateInBpm: 150, MaxHeartRateInBpm: 180, AvgCadenceInRpm: 90,
			}),
		}

		saver := NewExcelWorkoutPeriodSaverWithOptions(filePath, defaultSheet, "invalid")

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid start cell")
			})
		})
	})
}
