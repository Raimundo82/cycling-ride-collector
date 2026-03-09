package activity_excel

import (
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWorkoutExcelRecordMapperReturnsTheExpectedRowWithNormalWorkout(t *testing.T) {
	Convey("Given a normal workout and weight", t, func() {
		workout := domain.NewWorkout(domain.WorkoutParams{
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
			LegSensations:          "Boas",
		})

		Convey("When mapping", func() {
			row := workoutToRow(workout, 70.0)
			Convey("Then it returns the expected row values", func() {
				expected := []string{"6/1/2024", "Estrada", "10:30", "1h0m", "25.50", "500", "200", "220", "150", "180", "90", "Boas", "70.00"}
				So(row, ShouldResemble, expected)
			})
		})
	})
}

func TestWorkoutExcelRecordMapperReturnsTheExpectedRowWithRestWorkout(t *testing.T) {
	Convey("Given a rest workout and a weight", t, func() {
		workout := domain.NewWorkout(domain.WorkoutParams{
			WorkoutType:            domain.Descanso,
			StartTime:              time.Date(2024, 6, 1, 10, 30, 0, 0, time.UTC),
			DurationInMin:          -1,
			DistanceInKm:           -1,
			ElevationInMeters:      -1,
			AvgPowerInWatts:        -1,
			NormalizedPowerInWatts: -1,
			AvgHeartRateInBpm:      -1,
			MaxHeartRateInBpm:      -1,
			AvgCadenceInRpm:        -1,
			LegSensations:          "",
		})

		Convey("When mapping", func() {
			row := workoutToRow(workout, 70.0)
			Convey("Then it returns the expected row with empty values", func() {
				expected := []string{"6/1/2024", "Descanso", "", "", "", "", "", "", "", "", "", "", "70.00"}
				So(row, ShouldResemble, expected)
			})
		})
	})
}
