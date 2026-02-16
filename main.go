package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider"
	acitivityProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/provider/strava"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/repository/csv"
	authProvider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider/strava"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/repository/file"
	auth "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/service"
)

func main() {
	_ = godotenv.Load()
	cfg, request := parseFlagsAndConfig()

	if err := setupSaveWorkoutPeriodUseCase(cfg, request.DailyWorkoutPolicy).Execute(request.Period, request.MinimalWorkoutDuration); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Workout(s) processed and saved to %s (if any).\n", cfg.OutputFilePath)
	}
}

func parseFlagsAndConfig() (*config.Config, *input.SaveWorkoutPeriodRequest) {
	cfg := config.Load()

	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	outputFilePath := flag.String("output-file", "", "Output file path for the CSV")
	minimalWorkoutDuration := flag.Int("min-duration", 30, "Minimal workout duration in minutes")
	dailyWorkoutPolicy := flag.String("daily-workout-policy", "longest", "Daily workout policy")
	flag.Parse()

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

	period, err := domain.NewPeriod(startDate, endDate)
	if err != nil {
		log.Fatalf("Invalid period: %v", err)
	}

	if *dailyWorkoutPolicy != "longest" && *dailyWorkoutPolicy != "merge" {
		log.Fatalf("Invalid daily workout policy: %s. Allowed values are 'longest' or 'merge'", *dailyWorkoutPolicy)
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
		log.Fatalf("Invalid input: %v", err)
	}

	return cfg, request
}

func setupSaveWorkoutPeriodUseCase(cfg *config.Config, workoutPolicy string) usecase.SaveWorkoutPeriodUseCase {
	var dailyWorkoutPolicy contracts.DailyWorkoutPolicy

	switch workoutPolicy {
	case "longest":
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	case "merge":
		dailyWorkoutPolicy = usecase.NewMergeWorkouts()
	default:
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	}

	tokenRepo, err := file.NewTokenRepository(cfg.TokenFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize token repository: %v", err)
	}

	tokenProvider := auth.NewTokenService(authProvider.NewOAuthProvider(strava.NewOAuthHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg), cfg), tokenRepo)

	return usecase.NewSaveWorkoutPeriod(
		dailyWorkoutPolicy,
		csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath),
		provider.NewWorkoutProvider(acitivityProvider.NewActivityProvider(&http.Client{Timeout: 10 * time.Second}, cfg, tokenProvider)),
	)
}
