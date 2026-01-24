package csv

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCSVWorkoutRepository_SaveToWriter(t *testing.T) {
	Convey("Given a Workout", t, func() {
		workout := &domain.Workout{
			WorkoutType:            domain.Estrada,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          60,
			DistanceInKm:           25.5,
			ElevationInMeters:      500,
			AvgPowerInWatts:        200,
			NormalizedPowerInWatts: 220,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      180,
			AvgCadenceInRpm:        90,
		}
		repo := NewCSVWorkoutRepository("test_workouts.csv")

		Convey("When SaveToWriter is called", func() {
			var buf bytes.Buffer
			err := repo.SaveToWriter(workout, &buf)

			Convey("Then no error should occur", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the CSV output should be correct", func() {
				expected := "Estrada,10:30,60,25.50,500,200,220,150,180,90\n"
				So(buf.String(), ShouldEqual, expected)
			})
		})
	})
}

func TestCSVWorkoutRepository_Save_AppendsToFile(t *testing.T) {
	Convey("Given a file with initial content", t, func() {
		tmpfile, err := os.CreateTemp("", "workouts_temp.csv")
		So(err, ShouldBeNil)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		initialContent := "Estrada,10:00,30,10.00,100,150,160,120,130,80\n"
		_, err = tmpfile.WriteString(initialContent)
		So(err, ShouldBeNil)
		err = tmpfile.Close()
		So(err, ShouldBeNil)

		repo := NewCSVWorkoutRepository(tmpfile.Name())
		workout := &domain.Workout{
			WorkoutType:            domain.Estrada,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          60,
			DistanceInKm:           25.5,
			ElevationInMeters:      500,
			AvgPowerInWatts:        200,
			NormalizedPowerInWatts: 220,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      180,
			AvgCadenceInRpm:        90,
		}
		Convey("When Save is called", func() {
			err = repo.Save(workout)
			So(err, ShouldBeNil)

			data, err := os.ReadFile(tmpfile.Name())
			So(err, ShouldBeNil)
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			So(len(lines), ShouldEqual, 2)

			expectedFirst := "Estrada,10:00,30,10.00,100,150,160,120,130,80"
			expectedSecond := "Estrada,10:30,60,25.50,500,200,220,150,180,90"
			So(lines[0], ShouldEqual, expectedFirst)
			So(lines[1], ShouldEqual, expectedSecond)
		})
	})
}

func TestCSVWorkoutRepository_Save_CreatesFileIfNotExists(t *testing.T) {
	Convey("Given a non-existent file", t, func() {
		filename := "workouts_temp.csv"
		_ = os.Remove(filename)
		defer func() { _ = os.Remove(filename) }()
		repo := NewCSVWorkoutRepository(filename)
		workout := &domain.Workout{
			WorkoutType:            domain.Estrada,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          60,
			DistanceInKm:           25.5,
			ElevationInMeters:      500,
			AvgPowerInWatts:        200,
			NormalizedPowerInWatts: 220,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      180,
			AvgCadenceInRpm:        90,
		}
		Convey("When Save is called", func() {
			err := repo.Save(workout)
			So(err, ShouldBeNil)

			data, err := os.ReadFile(filename)
			So(err, ShouldBeNil)
			expected := "Estrada,10:30,60,25.50,500,200,220,150,180,90\n"
			So(string(data), ShouldEqual, expected)
		})
	})
}

func TestCSVWorkoutRepository_Save_InvalidFilePath(t *testing.T) {
	Convey("Given an invalid file path", t, func() {
		repo := NewCSVWorkoutRepository("/invalid/path/that/does/not/exist/workouts.csv")
		workout := &domain.Workout{
			WorkoutType:            domain.Estrada,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          60,
			DistanceInKm:           25.5,
			ElevationInMeters:      500,
			AvgPowerInWatts:        200,
			NormalizedPowerInWatts: 220,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      180,
			AvgCadenceInRpm:        90,
		}

		Convey("When Save is called", func() {
			err := repo.Save(workout)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCSVWorkoutRepository_Save_ReadOnlyFile(t *testing.T) {
	Convey("Given a read-only file", t, func() {
		tmpfile, err := os.CreateTemp("", "workouts_readonly.csv")
		So(err, ShouldBeNil)
		defer func() { _ = os.Remove(tmpfile.Name()) }()

		err = tmpfile.Close()
		So(err, ShouldBeNil)

		err = os.Chmod(tmpfile.Name(), 0o444)
		So(err, ShouldBeNil)

		repo := NewCSVWorkoutRepository(tmpfile.Name())
		workout := &domain.Workout{
			WorkoutType:            domain.Estrada,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          60,
			DistanceInKm:           25.5,
			ElevationInMeters:      500,
			AvgPowerInWatts:        200,
			NormalizedPowerInWatts: 220,
			AvgHeartRateInBpm:      150,
			MaxHeartRateInBpm:      180,
			AvgCadenceInRpm:        90,
		}

		Convey("When Save is called", func() {
			err := repo.Save(workout)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestCSVWorkoutRepository_SaveToWriter_NilWorkout(t *testing.T) {
	Convey("Given a nil workout", t, func() {
		repo := NewCSVWorkoutRepository("test_workouts.csv")

		Convey("When SaveToWriter is called with nil workout", func() {
			var buf bytes.Buffer

			Convey("Then it should panic", func() {
				So(func() { _ = repo.SaveToWriter(nil, &buf) }, ShouldPanic)
			})
		})
	})
}
