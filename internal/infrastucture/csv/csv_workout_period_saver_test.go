package csv

import (
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
