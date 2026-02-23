package activity_csv

import (
	"bytes"
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

var testAthlete = NewAthlete(70, 135, 240)

const testFilePath = "test_workouts.csv"

type errorOnWriteBuffer struct{}

func (b *errorOnWriteBuffer) Write(p []byte) (int, error) {
	return 0, os.ErrInvalid
}

func (b *errorOnWriteBuffer) Close() error {
	return nil
}

func TestCSVWorkoutPeriodSaverSavesWorkoutsToCSVFile(t *testing.T) {
	Convey("Given an array of workouts, an athlete and a CSV file", t, func() {
		tmpfile, _ := os.CreateTemp("", testFilePath)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

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

		csvPeriodSaver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := csvPeriodSaver.SaveAll(workouts, testAthlete)

			Convey("Then the CSV file should contain the correct data", func() {
				So(err, ShouldBeNil)
				data, _ := os.ReadFile(tmpfile.Name())
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")

				So(len(lines), ShouldEqual, 2)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Boas,70.00"
				expectedSecond := "6/2/2024,Estrada,09:00,0h45m,20.00,300,180,190,140,170,85,Médias,70.00"
				So(lines[0], ShouldEqual, expectedFirst)
				So(lines[1], ShouldEqual, expectedSecond)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverSavesSingleWorkoutToCSVFile(t *testing.T) {
	Convey("Given an array of workouts with single one, an athlete and a CSV file", t, func() {
		tmpfile, _ := os.CreateTemp("", testFilePath)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

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

		csvPeriodSaver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := csvPeriodSaver.SaveAll(workouts, testAthlete)

			Convey("Then the CSV file should contain the correct data", func() {
				So(err, ShouldBeNil)
				data, _ := os.ReadFile(tmpfile.Name())
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")

				So(len(lines), ShouldEqual, 1)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Boas,70.00"
				So(lines[0], ShouldEqual, expectedFirst)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverOverwritesWorkoutsToCSVFile(t *testing.T) {
	Convey("Given an array of workouts, an athlete and a file with a csv record", t, func() {
		tmpfile, _ := os.CreateTemp("", testFilePath)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		initialContent := "5/1/2024,Estrada,10:00,0h30m,10.00,100,150,160,120,130,80,Boas\n"
		_, _ = tmpfile.WriteString(initialContent)
		_ = tmpfile.Close()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

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
				LegSensations:          string(Medium),
			}),
		}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the CSV content should be replaced by the correct data", func() {
				So(err, ShouldBeNil)
				data, err := os.ReadFile(tmpfile.Name())
				So(err, ShouldBeNil)
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")
				So(len(lines), ShouldEqual, 1)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Médias,70.00"
				So(lines[0], ShouldEqual, expectedFirst)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverCreatesFileIfNotExist(t *testing.T) {
	Convey("Given a workout, an athlete and a non-existing CSV file", t, func() {
		defer func() { _ = os.Remove(testFilePath) }()
		saver := NewCSVWorkoutPeriodSaver(testFilePath)

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

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then the file should be created and contain the correct data", func() {
				So(err, ShouldBeNil)
				data, err := os.ReadFile(testFilePath)
				So(err, ShouldBeNil)
				expected := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Muito Boas,70.00\n"
				So(string(data), ShouldEqual, expected)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverReturnsErrorGivenAnInvalidFilePath(t *testing.T) {
	Convey("Given an invalid file path", t, func() {
		invalidPath := "/invalid/path/that/does/not/exist/workouts.csv"
		saver := NewCSVWorkoutPeriodSaver(invalidPath)
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

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "no such file or directory")
				So(err.Error(), ShouldContainSubstring, invalidPath)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverReturnsErrorGivenAReadOnlyFile(t *testing.T) {
	Convey("Given a read-only file", t, func() {
		tmpfile, _ := os.CreateTemp("", testFilePath)
		_ = tmpfile.Close()
		_ = os.Chmod(tmpfile.Name(), 0o444)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())
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
			}),
		}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "permission denied")
				So(err.Error(), ShouldContainSubstring, tmpfile.Name())
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverReturnsNoErrorGivenWorkoutWithZeroValues(t *testing.T) {
	Convey("Given a workout with only date and zero values for all other fields and an athlete", t, func() {
		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				ID:                     0,
				StartTime:              time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				WorkoutType:            Estrada,
				DistanceInKm:           0.00,
				DurationInMin:          0,
				ElevationInMeters:      0,
				AvgPowerInWatts:        0,
				NormalizedPowerInWatts: 0,
				AvgHeartRateInBpm:      0,
				MaxHeartRateInBpm:      0,
				AvgCadenceInRpm:        0,
				LegSensations:          string(Medium),
			}),
		}
		tmpfile, _ := os.CreateTemp("", testFilePath)
		defer func() { _ = os.Remove(tmpfile.Name()) }()
		_ = tmpfile.Close()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the file should contain a line with zeros", func() {
				data, err := os.ReadFile(tmpfile.Name())
				So(err, ShouldBeNil)
				expected := "6/1/2024,Estrada,00:00,0h0m,0.00,0,0,0,0,0,0,Médias,70.00\n"
				So(string(data), ShouldEqual, expected)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverReturnsNoErrorGivenWorkoutWithSentinelValues(t *testing.T) {
	Convey("Given a workout with only date and sentinel values for all other fields", t, func() {
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
		tmpfile, _ := os.CreateTemp("", testFilePath)
		defer func() { _ = os.Remove(tmpfile.Name()) }()
		_ = tmpfile.Close()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the file should contain an empty line", func() {
				data, err := os.ReadFile(tmpfile.Name())
				So(err, ShouldBeNil)
				expected := "6/1/2024,Descanso,,,,,,,,,,,70.00\n"
				So(string(data), ShouldEqual, expected)
			})
		})
	})
}

func TestCSVWorkoutPeriodSaverReturnsErrorWithInvalidClosingFileAction(t *testing.T) {
	Convey("Given a writer that returns an error on write", t, func() {
		errBuf := &errorOnWriteBuffer{}
		saver := &csvWorkoutPeriodSaver{
			filePath: "test_workouts.csv",
			buf:      &bytes.Buffer{},
			writer:   csv.NewWriter(errBuf),
		}
		workouts := []*Workout{
			NewWorkout(WorkoutParams{
				ID:                     0,
				StartTime:              time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				WorkoutType:            Estrada,
				DistanceInKm:           0.00,
				DurationInMin:          0,
				ElevationInMeters:      0,
				AvgPowerInWatts:        0,
				NormalizedPowerInWatts: 0,
				AvgHeartRateInBpm:      0,
				MaxHeartRateInBpm:      0,
				AvgCadenceInRpm:        0,
				LegSensations:          string(Medium),
			}),
		}

		Convey("When SaveAll is called", func() {
			err := saver.SaveAll(workouts, testAthlete)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid argument")
			})
		})
	})
}
