package models

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRideActivity_Valid(t *testing.T) {
	Convey("Given valid arguments", t, func() {
		dt := time.Now()
		duration := 5400

		Convey("When creating a ride activity", func() {
			activity, err := NewRideActivity(dt, duration)

			Convey("Then it stores the values and returns no error", func() {
				So(err, ShouldBeNil)
				So(activity.startDateTime, ShouldResemble, dt)
				So(activity.GetStartDateTime(), ShouldResemble, dt)
				So(activity.duration, ShouldEqual, duration)
				So(activity.GetDuration(), ShouldEqual, duration)
				So(activity.ToString(), ShouldEqual, fmt.Sprintf("StartDateTime: %s\nDuration: %d\n", dt, duration))
			})
		})
	})
}

func TestNewRideActivity_InvalidStartDateTime(t *testing.T) {
	Convey("Given invalid startDateTime", t, func() {
		dt := time.Time{}
		duration := 5400

		Convey("When creating a ride activity", func() {
			activity, err := NewRideActivity(dt, duration)

			Convey("Then it should not store the values and should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "invalid start date time")
				So(activity, ShouldBeNil)
			})
		})
	})
}

//func TestNewRideActivity_ZeroTime(t *testing.T) {
//	Convey("Given a zero time", t, func() {
//		dt := time.Time{}
//		_, err := NewRideActivity(dt, 5400)
//		So(err, ShouldNotBeNil)
//	})
//}
//
//func TestRideActivity_GetStartTime_Midnight(t *testing.T) {
//	Convey("Given a RideActivity at midnight", t, func() {
//		midnightDt := time.Date(2025, 12, 2, 0, 0, 0, 0, time.Local)
//		rideActivity, _ := NewRideActivity(midnightDt, 4500)
//		So(rideActivity.GetStartTime(), ShouldEqual, "00:00")
//	})
//}
//
//func TestRideActivity_GetStartTime_EndOfDay(t *testing.T) {
//	Convey("Given a RideActivity at end of day", t, func() {
//		rideActivity := RideActivity{
//			StartDateTime: time.Date(2025, 12, 2, 23, 59, 4, 4, time.Local),
//		}
//		So(rideActivity.GetStartTime(), ShouldEqual, "23:59")
//	})
//}
//
//func TestRideActivityGetStartDate(t *testing.T) {
//
//	Convey("Given a RideActivity with valid data", t, func() {
//		rideActivity := RideActivity{
//			StartDateTime: testDateTime,
//		}
//		Convey("When GetStartDate is called", func() {
//			result := rideActivity.GetStartDate()
//			Convey("Then the result should be the date part formatted", func() {
//				expected := "12/02/2025"
//				So(result, ShouldEqual, expected)
//			})
//		})
//	})
//}
//
//func TestRideActivityGetStartTime(t *testing.T) {
//
//	Convey("Given a RideActivity with valid data", t, func() {
//		rideActivity := RideActivity{
//			StartDateTime: testDateTime,
//		}
//		Convey("When GetStartDate is called", func() {
//			result := rideActivity.GetStartTime()
//			Convey("Then the result should be the time part formatted", func() {
//				expected := "08:30"
//				So(result, ShouldEqual, expected)
//			})
//		})
//	})
//}
//
//func TestRideActivityToString_Success(t *testing.T) {
//
//	Convey("Given a RideActivity with valid data", t, func() {
//		rideActivity := RideActivity{
//			StartDateTime: testDateTime,
//		}
//		Convey("When ToString is called", func() {
//			result := rideActivity.ToString()
//			Convey("Then the result should be the formatted date string", func() {
//				expected := "Date: 12/02/2025\nTime: 08:30\n"
//				So(result, ShouldEqual, expected)
//			})
//		})
//
//	})
//}
//
//func TestRideActivity_AlwaysStoresInLisbonTimeZone(t *testing.T) {
//	Convey("Given a start date time in New York time zone", t, func() {
//		nyLoc, err := time.LoadLocation("America/New_York")
//		So(err, ShouldBeNil)
//		nyTime := testDateTime.In(nyLoc)
//
//		Convey("When creating a RideActivity", func() {
//			So(err, ShouldBeNil)
//			activity, _ := NewRideActivity(nyTime, 4500)
//
//			Convey("Then StartDateTime should be in Europe/Lisbon time", func() {
//				So(activity.StartDateTime.Location().String(), ShouldEqual, "Europe/Lisbon")
//				expectedLisbonTime := testDateTime
//				So(activity.StartDateTime.Hour(), ShouldEqual, expectedLisbonTime.Hour())
//				So(activity.StartDateTime.Minute(), ShouldEqual, expectedLisbonTime.Minute())
//			})
//		})
//	})
//}
//
