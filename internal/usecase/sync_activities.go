package usecase

import (
	"fmt"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

// SyncActivitiesUseCase implements the business logic for syncing activities
type SyncActivitiesUseCase struct {
	activitiesRepo domain.ActivitiesRepository
	sheetsWriter   domain.SheetsWriter
}

// NewSyncActivitiesUseCase creates a new instance of SyncActivitiesUseCase
func NewSyncActivitiesUseCase(
	activitiesRepo domain.ActivitiesRepository,
	sheetsWriter domain.SheetsWriter,
) *SyncActivitiesUseCase {
	return &SyncActivitiesUseCase{
		activitiesRepo: activitiesRepo,
		sheetsWriter:   sheetsWriter,
	}
}

// SyncWeeklyActivities fetches activities from the last week and writes them to sheets
func (s *SyncActivitiesUseCase) SyncWeeklyActivities() error {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)

	fmt.Printf("Fetching activities from %s to %s\n",
		weekAgo.Format(time.RFC3339),
		now.Format(time.RFC3339))

	activities, err := s.activitiesRepo.GetActivitiesBetween(weekAgo, now)
	if err != nil {
		return fmt.Errorf("failed to fetch activities: %w", err)
	}

	fmt.Printf("Found %d activities\n", len(activities))

	if len(activities) == 0 {
		fmt.Println("No activities to sync")
		return nil
	}

	if err := s.sheetsWriter.WriteActivities(activities); err != nil {
		return fmt.Errorf("failed to write activities to sheets: %w", err)
	}

	fmt.Println("Successfully synced activities to Google Sheets")
	return nil
}
