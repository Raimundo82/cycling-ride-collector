package models

import (
	"fmt"
	"time"
)

type RideActivity struct {
	startDateTime time.Time
	duration      int
}

func NewRideActivity(startDateTime time.Time, duration int) (*RideActivity, error) {
	acitivity := &RideActivity{}
	if err := acitivity.SetStartDateTime(startDateTime); err != nil {
		return nil, err
	}
	acitivity.SetDuration(duration)
	return acitivity, nil
}

func (activity RideActivity) GetStartDateTime() time.Time {
	return activity.startDateTime
}

func (activity *RideActivity) SetStartDateTime(startDateTime time.Time) error {
	if startDateTime.IsZero() {
		return fmt.Errorf("invalid start date time")
	}
	activity.startDateTime = startDateTime
	return nil
}

func (activity RideActivity) GetDuration() int {
	return activity.duration
}

func (activity *RideActivity) SetDuration(duration int) {
	activity.duration = duration
}

//	func (activity RideActivity) GetStartDate() string {
//		return activity.startDateTime.Format("01/02/2006")
//	}
//
//	func (activity RideActivity) GetStartTime() string {
//		return activity.StartDateTime.Format("15:04")
//	}
//
//	func (activity RideActivity) GetDuration() string {
//		hours := activity.Duration / 3600
//		minutes := (activity.Duration % 3600) % 60
//		return fmt.Sprintf("%dh%02d'", hours, minutes)
//	}
func (activity RideActivity) ToString() string {
	return fmt.Sprintf("StartDateTime: %s\nDuration: %d\n", activity.GetStartDateTime().String(), activity.GetDuration())
}
