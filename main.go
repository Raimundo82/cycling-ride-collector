package main

import (
	"fmt"
	"os"

	"github.com/raimundo82/go-strava-weekly/internal/config"
	"github.com/raimundo82/go-strava-weekly/internal/infrastructure/sheets"
	"github.com/raimundo82/go-strava-weekly/internal/infrastructure/strava"
	"github.com/raimundo82/go-strava-weekly/internal/interfaces/cli"
	"github.com/raimundo82/go-strava-weekly/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize infrastructure layer (adapters)
	stravaClient := strava.NewClient(cfg.StravaAPIKey)
	sheetsClient := sheets.NewClient(cfg.SheetsSpreadsheetID)

	// Initialize use case layer (business logic)
	syncUseCase := usecase.NewSyncActivitiesUseCase(stravaClient, sheetsClient)

	// Initialize interface layer (CLI)
	app := cli.NewApp(syncUseCase)

	// Run the application
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
