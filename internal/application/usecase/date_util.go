package usecase

import (
	"fmt"
	"time"
)

func ParseDateToUnix(dateStr string) (int64, error) {
	if dateStr == "" {
		return 0, fmt.Errorf("empty date string")
	}

	const layout = "2006-01-02"
	t, err := time.ParseInLocation(layout, dateStr, time.UTC)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %w", err)
	}

	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return midnight.Unix(), nil
}
