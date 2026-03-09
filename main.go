package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	log.Println("Workout summary saved successfully.")
}

func run() error {
	_ = godotenv.Load()
	cfg, request, err := parseFlagsAndConfig()
	if err != nil {
		return err
	}

	app, err := NewApp(cfg, request.DailyWorkoutPolicy)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	if err := app.Run(request); err != nil {
		return fmt.Errorf("error running app: %w", err)
	}

	return nil
}

func parseFlagsAndConfig() (*config.Config, *input.SaveWorkoutPeriodRequest, error) {
	cfg := config.Load()

	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	outputFilePath := flag.String("output-file", "", "Output file path for the CSV")
	minimalWorkoutDuration := flag.Int("min-duration", 30, "Minimal workout duration in minutes")
	dailyWorkoutPolicy := flag.String("daily-workout-policy", "longest", "Daily workout policy")
	flag.Parse()

	if *startDateStr == "" || *endDateStr == "" {
		return nil, nil, errors.New("flags --start-date and --end-date are required")
	}

	const layout = "01/02/2006"
	startDate, err := time.Parse(layout, *startDateStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse(layout, *endDateStr)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid end date: %w", err)
	}

	period, err := domain.NewPeriod(startDate, endDate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid period: %w", err)
	}

	if *dailyWorkoutPolicy != "longest" && *dailyWorkoutPolicy != "merge" {
		return nil, nil, fmt.Errorf("invalid daily workout policy: %s. Allowed values are 'longest' or 'merge'", *dailyWorkoutPolicy)
	}

	if *minimalWorkoutDuration < 30 {
		*minimalWorkoutDuration = 30
	}

	if *outputFilePath == "" {
		*outputFilePath = fmt.Sprintf("workouts_summary_%s_to_%s.csv", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}
	cfg.OutputFilePath = *outputFilePath

	request, err := input.NewSaveWorkoutPeriodRequest(period, *dailyWorkoutPolicy, *minimalWorkoutDuration)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid input: %w", err)
	}

	return cfg, request, nil
}
