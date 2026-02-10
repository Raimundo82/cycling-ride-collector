package domain

import (
	"errors"
	"time"
)

type Period interface {
	StartDate() time.Time
	EndDate() time.Time
}

type period struct {
	startDate time.Time
	endDate   time.Time
}

func NewPeriod(startDate, endDate time.Time) (Period, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, errors.New("start date and end date must be valid")
	}

	if startDate.After(endDate) || startDate.Equal(endDate) {
		return nil, errors.New("start date must be before end date")
	}
	return &period{
		startDate: startDate,
		endDate:   endDate,
	}, nil
}

func (p *period) StartDate() time.Time {
	return p.startDate
}

func (p *period) EndDate() time.Time {
	return p.endDate
}
