package domain

import "time"

// Activity represents a Strava activity entity
type Activity struct {
	ID             int64
	Name           string
	Type           string
	Distance       float64 // meters
	MovingTime     int     // seconds
	ElapsedTime    int     // seconds
	TotalElevation float64 // meters
	StartDate      time.Time
	AverageSpeed   float64 // meters per second
}

// ActivitiesRepository defines the interface for fetching activities
type ActivitiesRepository interface {
	GetActivitiesBetween(startDate, endDate time.Time) ([]Activity, error)
}

// SheetsWriter defines the interface for writing activities to sheets
type SheetsWriter interface {
	WriteActivities(activities []Activity) error
}
