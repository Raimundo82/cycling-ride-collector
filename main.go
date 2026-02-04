package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/go-strava-weekly/internal/application/contracts"
	"github.com/raimundo82/go-strava-weekly/internal/application/orchestration"
	"github.com/raimundo82/go-strava-weekly/internal/application/usecase"
	"github.com/raimundo82/go-strava-weekly/internal/config"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/csv"
	"github.com/raimundo82/go-strava-weekly/internal/infrastucture/strava"
)

func main() {
	_ = godotenv.Load()
	cfg := parseFlagsAndConfig()

	orchestrator := setupOrchestrator(cfg)
	period := orchestration.NewPeriod(cfg.StartDate, cfg.EndDate)

	if err := orchestrator.SaveWorkoutsOverPeriod(period, cfg.MinimalWorkoutDuration); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Workout(s) processed and saved to %s (if any).\n", cfg.OutputFilePath)
	}
}

func parseFlagsAndConfig() *config.Config {
	cfg := config.Load()

	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	accessToken := flag.String("access-token", "", "Strava API access token")
	outputFilePath := flag.String("output-file", "", "Output file path for the CSV")
	minimalWorkoutDuration := flag.Int("min-duration", 0, "Minimal workout duration in minutes")
	dailyWorkoutPolicy := flag.String("daily-workout-policy", "", "Daily workout policy")
	flag.Parse()

	if *accessToken != "" {
		cfg.StravaAccessToken = *accessToken
	}
	if *dailyWorkoutPolicy != "" {
		cfg.DailyWorkoutPolicy = *dailyWorkoutPolicy
	}

	if *minimalWorkoutDuration > 0 {
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
	cfg.StartDate = startDate

	endDate, err := time.Parse(layout, *endDateStr)
	if err != nil {
		log.Fatalf("Invalid end date: %v", err)
	}
	cfg.EndDate = endDate

	if *outputFilePath != "" {
		cfg.OutputFilePath = *outputFilePath
	} else if cfg.OutputFilePath == "" {
		cfg.OutputFilePath = fmt.Sprintf("workouts_summary_%s_to_%s.csv", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}
	return cfg
}

func setupOrchestrator(cfg *config.Config) *orchestration.SaveWorkoutsOrchestrator {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	stravaClient := strava.NewHttpClient(httpClient, cfg)

	var dailyWorkoutPolicy contracts.DailyWorkoutPolicy
	switch cfg.DailyWorkoutPolicy {
	case "longest":
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	case "merge":
		dailyWorkoutPolicy = usecase.NewMergeWorkouts()
	default:
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	}

	save := usecase.NewSaveWorkout(
		dailyWorkoutPolicy,
		csv.NewCSVWorkoutRepository(cfg.OutputFilePath),
		strava.NewProvider(stravaClient),
	)
	return &orchestration.SaveWorkoutsOrchestrator{
		SaveWorkoutUseCase: save,
	}
}
