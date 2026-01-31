package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/go-strava-weekly/internal/application/orchestration"
	"github.com/raimundo82/go-strava-weekly/internal/application/usecase"
	"github.com/raimundo82/go-strava-weekly/internal/config"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/csv"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/strava"
)

func main() {
	_ = godotenv.Load()
	cfg, startDate, endDate, outputFilePath := parseFlagsAndConfig()

	orchestrator := setupOrchestrator(cfg, outputFilePath)
	period := orchestration.NewPeriod(startDate, endDate)

	if err := orchestrator.SaveWorkoutsOverPeriod(period, cfg.MinimalWorkoutDuration); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Workout(s) processed and saved to %s (if any).\n", outputFilePath)
	}
}

func parseFlagsAndConfig() (*config.Config, time.Time, time.Time, string) {
	cfg := config.Load()

	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	accessToken := flag.String("access-token", "", "Strava API access token")
	outputFilePath := flag.String("output-file", "", "Output file path for the CSV")
	minimalWorkoutDuration := flag.Int("min-duration", cfg.MinimalWorkoutDuration, "Minimal workout duration in minutes")
	flag.Parse()

	if *accessToken != "" {
		cfg.StravaAccessToken = *accessToken
	}
	if flag.Lookup("min-duration").Value.String() != fmt.Sprint(cfg.MinimalWorkoutDuration) {
		cfg.MinimalWorkoutDuration = *minimalWorkoutDuration
	}

	if *startDateStr == "" || *endDateStr == "" {
		log.Fatal("Flags --start-date and --end-date are required")
	}
	const layout = "01/02/2006"
	startDate, err := time.Parse(layout, *startDateStr)
	if err != nil {
		log.Fatalf("Invalid start date: %v", err)
	}
	endDate, err := time.Parse(layout, *endDateStr)
	if err != nil {
		log.Fatalf("Invalid end date: %v", err)
	}

	if *outputFilePath == "" {
		*outputFilePath = fmt.Sprintf("workouts_summary_%s_to_%s.csv", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	return cfg, startDate, endDate, *outputFilePath
}

func setupOrchestrator(cfg *config.Config, outputFilePath string) *orchestration.SaveWorkoutsOrchestrator {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	stravaClient := strava.NewHttpClient(httpClient, cfg)

	save := &usecase.SaveWorkout{
		WorkoutRepo:     csv.NewCSVWorkoutRepository(outputFilePath),
		WorkoutProvider: strava.NewProvider(stravaClient),
	}

	return &orchestration.SaveWorkoutsOrchestrator{
		SaveWorkoutUseCase: save,
	}
}
