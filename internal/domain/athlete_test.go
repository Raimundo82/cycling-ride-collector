package domain

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewAthleteReturnsValidAthlete(t *testing.T) {
	Convey("Given valid athlete parameters", t, func() {
		weightInKg := 70.0
		heartRateThresholdInBpm := 135
		powerThresholdInWatts := 250

		Convey("When NewAthlete is called", func() {
			athlete := NewAthlete(weightInKg, heartRateThresholdInBpm, powerThresholdInWatts)

			Convey("Then the returned athlete should have the correct attributes", func() {
				So(athlete.WeightInKg(), ShouldEqual, weightInKg)
				So(athlete.HeartRateThresholdInBpm(), ShouldEqual, heartRateThresholdInBpm)
				So(athlete.PowerThresholdInWatts(), ShouldEqual, powerThresholdInWatts)
			})
		})
	})
}

func TestNewAthleteReturnsValidAthleteWithZeroValues(t *testing.T) {
	Convey("Given valid athlete parameters", t, func() {
		weightInKg := -70.0
		heartRateThresholdInBpm := -135
		powerThresholdInWatts := -250

		Convey("When NewAthlete is called", func() {
			athlete := NewAthlete(weightInKg, heartRateThresholdInBpm, powerThresholdInWatts)

			Convey("Then the returned athlete should have the correct attributes", func() {
				So(athlete.WeightInKg(), ShouldEqual, 0)
				So(athlete.HeartRateThresholdInBpm(), ShouldEqual, 0)
				So(athlete.PowerThresholdInWatts(), ShouldEqual, 0)
			})
		})
	})
}
