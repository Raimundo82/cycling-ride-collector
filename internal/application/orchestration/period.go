package orchestration

import "time"

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func NewDate(year int, month time.Month, day int) Date {
	return Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

type Period struct {
	StartDate time.Time
	EndDate   time.Time
}

func NewPeriod(startDate, endDate Date) Period {
	return Period{
		StartDate: startDate.ToTime(),
		EndDate:   endDate.ToTime(),
	}
}
