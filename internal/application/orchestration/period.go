package orchestration

import (
	"time"
)

type Period struct {
	StartDate time.Time
	EndDate   time.Time
}

func NewPeriod(startDate, endDate time.Time) Period {
	return Period{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func (p Period) IsValid() bool {
	return !p.StartDate.After(p.EndDate)
}
