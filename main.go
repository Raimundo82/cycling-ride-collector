package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase/input"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/csv"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/strava"
	"gofr.dev/pkg/gofr"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	if isCronMode() {
		runCron(cfg)
		return
	}

	runCLI(cfg)
}

func runCLI(cfg *config.Config) {
	cfg, request := parseFlagsAndConfig(cfg)

	uc := setupSaveWorkoutPeriodUseCase(cfg, request.DailyWorkoutPolicy)
	executeAndReport(uc, request.Period, request.MinimalWorkoutDuration, cfg.OutputFilePath)
}

func runCron(cfg *config.Config) {
	app := gofr.New()
	app.AddCronJob("0 0 19 * * 0", "weekly_sunday_7pm", func(ctx *gofr.Context) {
		endDate := time.Now()
		startDate := endDate.Add(-6 * 24 * time.Hour)
		ctx.Logger.Infof("Processing workouts from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

		period, err := domain.NewPeriod(startDate, endDate)
		if err != nil {
			ctx.Logger.Errorf("invalid period: %v", err)
			return
		}

		cfg.OutputFilePath = generateOutputFilePath(startDate, endDate)

		uc := usecase.NewSaveWorkoutPeriod(
			usecase.NewLongestWorkout(),
			csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath),
			strava.NewProvider(strava.NewHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg)))

		executeAndReport(uc, period, 30, cfg.OutputFilePath)
	})
}

func executeAndReport(uc usecase.SaveWorkoutPeriodUseCase, period domain.Period, minDuration int, outputPath string) {
	if err := uc.Execute(period, minDuration); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Workout(s) processed and saved to %s (if any).\n", outputPath)
	}
}

func parseFlagsAndConfig(cfg *config.Config) (*config.Config, *input.SaveWorkoutPeriodRequest) {
	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	accessToken := flag.String("access-token", "", "Strava API access token")
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
		*outputFilePath = generateOutputFilePath(startDate, endDate)
	}
	cfg.OutputFilePath = *outputFilePath

	if *accessToken != "" {
		cfg.StravaAccessToken = *accessToken
	}

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

	return usecase.NewSaveWorkoutPeriod(
		dailyWorkoutPolicy,
		csv.NewCSVWorkoutPeriodSaver(cfg.OutputFilePath),
		strava.NewProvider(strava.NewHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg)),
	)
}

func isCronMode() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--cron" || arg == "-cron" {
			return true
		}
	}
	return false
}

func generateOutputFilePath(startDate, endDate time.Time) string {
	return fmt.Sprintf("workouts_summary_%s_to_%s.csv", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
}
