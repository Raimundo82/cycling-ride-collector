package csv

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCsvWorkoutPeriodSaver_GivenAnArrayOfWorkoutsAndACSVFile_WhenSaveAllIsCalled_ThenTheCSVFileShouldContainTheCorrectData(t *testing.T) {
	Convey("Given an array of workouts and a CSV file", t, func() {
		tmpfile, _ := os.CreateTemp("", "workouts_temp.csv")
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Boas",
			}),
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Médias",
			}),
		}

		csvPeriodSaver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := csvPeriodSaver.SaveAll(workouts)

			Convey("Then the CSV file should contain the correct data", func() {
				So(err, ShouldBeNil)
				data, _ := os.ReadFile(tmpfile.Name())
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")

				So(len(lines), ShouldEqual, 2)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Boas"
				expectedSecond := "6/2/2024,Estrada,09:00,0h45m,20.00,300,180,190,140,170,85,Médias"
				So(lines[0], ShouldEqual, expectedFirst)
				So(lines[1], ShouldEqual, expectedSecond)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenAnArrayOfWorkoutsWithSingleOneAndACSVFile_WhenSaveAllIsCalled_ThenTheCSVFileShouldContainTheCorrectData(t *testing.T) {
	Convey("Given an array of workouts with single one and a CSV file", t, func() {
		tmpfile, _ := os.CreateTemp("", "workouts_temp.csv")
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Boas",
			}),
		}

		csvPeriodSaver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		Convey("When SaveAll is called", func() {
			err := csvPeriodSaver.SaveAll(workouts)

			Convey("Then the CSV file should contain the correct data", func() {
				So(err, ShouldBeNil)
				data, _ := os.ReadFile(tmpfile.Name())
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")

				So(len(lines), ShouldEqual, 1)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Boas"
				So(lines[0], ShouldEqual, expectedFirst)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenAnArrayOfWorkoutsAndACSVFileWithCsvRecords_WhenSaveAllIsCalled_ThenTheCSVFileContentShouldBeReplacedByTheCorrectData(t *testing.T) {
	Convey("Given an array of workouts and a file with a csv record", t, func() {
		tmpfile, _ := os.CreateTemp("", "workouts_temp.csv")
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		initialContent := "5/1/2024,Estrada,10:00,0h30m,10.00,100,150,160,120,130,80,Boas\n"
		_, _ = tmpfile.WriteString(initialContent)
		_ = tmpfile.Close()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())

		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Médias",
			}),
		}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts)

			Convey("Then the CSV content should be replaced by the correct data", func() {
				So(err, ShouldBeNil)
				data, err := os.ReadFile(tmpfile.Name())
				So(err, ShouldBeNil)
				lines := strings.Split(strings.TrimSpace(string(data)), "\n")
				So(len(lines), ShouldEqual, 1)
				expectedFirst := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Médias"
				So(lines[0], ShouldEqual, expectedFirst)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenANonExistentFile_WhenSaveIsCalled_ThenTheFileIsCreatedAndContainsTheCorrectData(t *testing.T) {
	Convey("Given a non-existent file", t, func() {
		filename := "workouts_temp.csv"
		defer func() { _ = os.Remove(filename) }()
		saver := NewCSVWorkoutPeriodSaver(filename)

		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Muito Boas",
			}),
		}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts)

			Convey("Then the file should be created and contain the correct data", func() {
				So(err, ShouldBeNil)
				data, err := os.ReadFile(filename)
				So(err, ShouldBeNil)
				expected := "6/1/2024,Estrada,10:30,1h0m,25.50,500,200,220,150,180,90,Muito Boas\n"
				So(string(data), ShouldEqual, expected)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenAnInvalidFilePath_WhenSaveIsCalled_ThenAnErrorIsReturned(t *testing.T) {
	Convey("Given an invalid file path", t, func() {
		saver := NewCSVWorkoutPeriodSaver("/invalid/path/that/does/not/exist/workouts.csv")
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Muito Boas",
			}),
		}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenAReadOnlyFile_WhenSaveIsCalled_ThenAnErrorIsReturned(t *testing.T) {
	Convey("Given a read-only file", t, func() {
		tmpfile, _ := os.CreateTemp("", "workouts_readonly.csv")
		_ = tmpfile.Close()
		_ = os.Chmod(tmpfile.Name(), 0o444)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		saver := NewCSVWorkoutPeriodSaver(tmpfile.Name())
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
			})}

		Convey("When Save is called", func() {
			err := saver.SaveAll(workouts)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenWorkoutWithZeroValues_WhenSaveIsCalled_ThenNoErrorIsReturned(t *testing.T) {
	Convey("Given a workout with only date and zero values for all other fields", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
				LegSensations:          "Médias",
			})}
		saver := &csvWorkoutPeriodSaver{filePath: "test_workouts.csv"}

		Convey("When SaveToWriter is called", func() {
			var buf bytes.Buffer
			err := saver.SaveToWriterAll(workouts, &buf)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("And it should write a line with zeros", func() {
				expected := "6/1/2024,Estrada,00:00,0h0m,0.00,0,0,0,0,0,0,Médias\n"
				So(buf.String(), ShouldEqual, expected)
			})
		})
	})
}

func TestCsvWorkoutPeriodSaver_GivenWorkoutWithSentinelValues_WhenSaveIsCalled_ThenNoErrorIsReturned(t *testing.T) {
	Convey("Given a workout with only date and sentinel values for all other fields", t, func() {
		workouts := []*Workout{
			NewWorkout(&WorkoutParams{
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
			})}
		saver := &csvWorkoutPeriodSaver{filePath: "test_workouts.csv"}

		Convey("When SaveToWriter is called", func() {
			var buf bytes.Buffer
			err := saver.SaveToWriterAll(workouts, &buf)

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("And it should write an empty line", func() {
				expected := "6/1/2024,Descanso,,,,,,,,,,\n"
				So(buf.String(), ShouldEqual, expected)
			})
		})
	})
}
