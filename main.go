package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/scheduler"
)

const dateLayout = "01/02/2006"

var startWeeklySunday20 = func(job func()) error {
	s := scheduler.New()
	return s.StartWeeklySunday20(job)
}

var waitForCronMode = func() {
	select {}
}

func main() {
	_ = godotenv.Load()
	options, err := config.ParseCLI(os.Args[1:])
	if err != nil {
		log.Printf("Error parsing CLI: %v", err)
		os.Exit(1)
	}

	if err := run(options); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	if options.CronMode {
		log.Println("Cron scheduler started. Waiting for weekly runs...")
		waitForCronMode()
		return
	}
	log.Println("Workout summary saved successfully.")
}

func run(options config.CLIOptions) error {
	if options.CronMode {
		return runCronMode()
	}
	return runOnceMode(options)
}

func runCronMode() error {
	return startWeeklySunday20(func() {
		cronOptions := buildCronOptions(time.Now().UTC())

		if err := runOnceMode(cronOptions); err != nil {
			log.Printf("scheduled run failed: %v", err)
		}
	})
}

func buildCronOptions(now time.Time) config.CLIOptions {
	start := now.AddDate(0, 0, -6).Format(dateLayout)
	end := now.Format(dateLayout)

	return config.CLIOptions{
		CronMode:               true,
		StartDate:              start,
		EndDate:                end,
		DailyWorkoutPolicy:     "longest",
		MinimalWorkoutDuration: 30,
	}
}

func runOnceMode(options config.CLIOptions) error {
	cfg := buildRuntimeConfig(options)
	if err := cfg.ValidateRequired(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	request, err := buildSaveWorkoutRequest(options)
	if err != nil {
		return err
	}
	if cfg.OutputFilePath == "" {
		cfg.OutputFilePath = fmt.Sprintf(
			"workouts_summary_%s_to_%s.xlsx",
			request.Period.StartDate().Format("2006-01-02"),
			request.Period.EndDate().Format("2006-01-02"),
		)
	}

	app, err := NewApp(cfg, request.DailyWorkoutPolicy)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	if err := app.Run(request, cfg.OutputFilePath); err != nil {
		return fmt.Errorf("error running app: %w", err)
	}

	return nil
}

func buildRuntimeConfig(options config.CLIOptions) *config.Config {
	cfg := config.Load()
	if options.OutputFilePath != "" {
		cfg.OutputFilePath = options.OutputFilePath
	}
	return cfg
}

func buildSaveWorkoutRequest(options config.CLIOptions) (*input.SaveWorkoutPeriodRequest, error) {
	if options.StartDate == "" || options.EndDate == "" {
		return nil, errors.New("flags --start-date and --end-date are required")
	}

	startDate, err := time.Parse(dateLayout, options.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse(dateLayout, options.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	period, err := domain.NewPeriod(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}

	if options.DailyWorkoutPolicy != "longest" && options.DailyWorkoutPolicy != "merge" {
		return nil, fmt.Errorf("invalid daily workout policy: %s. Allowed values are 'longest' or 'merge'", options.DailyWorkoutPolicy)
	}

	minimalDuration := options.MinimalWorkoutDuration
	if minimalDuration < 30 {
		minimalDuration = 30
	}

	request, err := input.NewSaveWorkoutPeriodRequest(period, options.DailyWorkoutPolicy, minimalDuration)
	if err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	return request, nil
}
