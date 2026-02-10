package domain

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPeriod_GivenEndDateDayAfterStartDate_WhenInvoked_ThenValidPeriodReturned(t *testing.T) {
	Convey("Given a start date and an end date that is the day after", t, func() {
		startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)

		Convey("When NewPeriod is invoked", func() {
			p, err := NewPeriod(startDate, endDate)

			Convey("Then a valid Period is returned", func() {
				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.StartDate().Equal(startDate), ShouldBeTrue)
				So(p.EndDate().Equal(endDate), ShouldBeTrue)
			})
		})
	})
}

func TestNewPeriod_GivenInvalidStartDate_WhenInvoked_ThenErrorIsReturned(t *testing.T) {
	Convey("Given an invalid start date", t, func() {
		startDate := time.Time{}
		endDate := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)

		Convey("When NewPeriod is invoked", func() {
			p, err := NewPeriod(startDate, endDate)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(p, ShouldBeNil)
			})
		})
	})
}

func TestNewPeriod_GivenInvalidEndDate_WhenInvoked_ThenErrorIsReturned(t *testing.T) {
	Convey("Given an invalid end date", t, func() {
		startDate := time.Time{}
		endDate := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)

		Convey("When NewPeriod is invoked", func() {
			p, err := NewPeriod(startDate, endDate)

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(p, ShouldBeNil)
			})
		})
	})
}

func TestNewPeriod_GivenEndDateInTheSameDayAsStartDate_WhenInvoked_ThenAOneDayPeriodIsReturned(t *testing.T) {
	Convey("Given a start date and an end date that is the same day", t, func() {
		startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

		Convey("When NewPeriod is invoked", func() {
			p, err := NewPeriod(startDate, endDate)

			Convey("Then a valid Period is returned", func() {
				So(err, ShouldBeNil)
				So(p, ShouldNotBeNil)
				So(p.StartDate().Equal(startDate), ShouldBeTrue)
				So(p.EndDate().Equal(endDate), ShouldBeTrue)
			})
		})
	})
}

func TestGetDatesImmutability_GivenNewPeriod_WhenGettingStartOrEndDate_ThenOriginalDatesShouldNotBeModified(t *testing.T) {
	Convey("Given a new Period", t, func() {
		startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC)
		p, _ := NewPeriod(startDate, endDate)

		Convey("When getting and mutating the start date ", func() {
			s := p.StartDate()
			s = s.Add(24 * time.Hour)

			Convey("Then the original start date should not be modified", func() {
				So(p.StartDate().Equal(s), ShouldBeFalse)
				So(p.StartDate().Equal(startDate), ShouldBeTrue)
			})
		})

		Convey("When getting and mutating the end date ", func() {
			e := p.EndDate()
			e = e.Add(24 * time.Hour)

			Convey("Then the original end date should not be modified", func() {
				So(p.EndDate().Equal(e), ShouldBeFalse)
				So(p.EndDate().Equal(endDate), ShouldBeTrue)
			})
		})
	})
}
