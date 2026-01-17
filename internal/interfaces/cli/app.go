package cli

import (
	"fmt"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

// App represents the CLI application
type App struct {
	syncService domain.SyncService
}

// NewApp creates a new CLI application
func NewApp(syncService domain.SyncService) *App {
	return &App{
		syncService: syncService,
	}
}

// Run executes the application
func (a *App) Run() error {
	fmt.Println("== Strava → Google Sheets Sync ==")

	if err := a.syncService.SyncWeeklyActivities(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return nil
}
