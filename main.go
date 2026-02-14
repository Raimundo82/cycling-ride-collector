package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/usecase"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/csv"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastucture/strava"
	"gofr.dev/pkg/gofr"
)

func main() {
	_ = godotenv.Load()
	flags := getFlags()
	cfg := config.Load()

	if flags.cronMode {
		runCron(cfg, flags)
		return
	}

	err := runCLI(cfg, flags)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(cfg *config.Config, flags *flags) error {
	accessToken := flags.accessToken

	period, err := parseDateRange(flags.startDate, flags.endDate)
	if err != nil {
		return err
	}

	dailyWorkoutPolicy, err := validateWorkoutPolicy(flags.dailyWorkoutPolicy)
	if err != nil {
		return err
	}

	minimalWorkoutDuration := getMinWorkoutDuration(flags.minimalWorkoutDuration)

	outputFilePath := flags.outputFilePath
	if outputFilePath == "" {
		outputFilePath = generateOutputFilePath(*period)
	}

	if accessToken != "" {
		cfg.StravaAccessToken = accessToken
	}

	executeUsecase(cfg, period, dailyWorkoutPolicy, minimalWorkoutDuration, outputFilePath)
	return nil
}

func runCron(cfg *config.Config, flags *flags) {
	app := gofr.New()
	app.AddCronJob("0 19 * * 0", "weekly_sunday_7pm", func(ctx *gofr.Context) {
		endDate := time.Now()
		startDate := endDate.AddDate(0, 0, -6)
		ctx.Logger.Infof("Processing workouts from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

		period, err := domain.NewPeriod(startDate, endDate)
		if err != nil {
			ctx.Logger.Errorf("invalid period: %v", err)
			return
		}

		outputFilePath := generateOutputFilePath(period)
		minimalWorkoutDuration := getMinWorkoutDuration(flags.minimalWorkoutDuration)
		executeUsecase(cfg, &period, flags.dailyWorkoutPolicy, minimalWorkoutDuration, outputFilePath)
	})
	app.Run()
}

func executeUsecase(cfg *config.Config, period *domain.Period, workoutPolicy string, minimalWorkoutDuration int, outputPath string) {
	var dailyWorkoutPolicy contracts.DailyWorkoutPolicy

	switch workoutPolicy {
	case "longest":
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	case "merge":
		dailyWorkoutPolicy = usecase.NewMergeWorkouts()
	default:
		dailyWorkoutPolicy = usecase.NewLongestWorkout()
	}

	uc := usecase.NewSaveWorkoutPeriod(
		dailyWorkoutPolicy,
		csv.NewCSVWorkoutPeriodSaver(outputPath),
		strava.NewProvider(strava.NewHttpClient(&http.Client{Timeout: 10 * time.Second}, cfg)),
	)

	if err := uc.Execute(*period, minimalWorkoutDuration); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Workout(s) processed and saved to %s (if any).\n", outputPath)
	}
}

func generateOutputFilePath(period domain.Period) string {
	return fmt.Sprintf("workouts_summary_%s_to_%s.csv", period.StartDate().Format("2006-01-02"), period.EndDate().Format("2006-01-02"))
}

func getMinWorkoutDuration(duration int) int {
	if duration < 30 {
		return 30
	}
	return duration
}

func parseDateRange(startStr, endStr string) (*domain.Period, error) {
	const layout = "01/02/2006"
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse(layout, endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}
	period, err := domain.NewPeriod(start, end)
	if err != nil {
		return nil, fmt.Errorf("invalid period: %w", err)
	}
	return &period, nil
}

func validateWorkoutPolicy(policy string) (string, error) {
	if policy != "longest" && policy != "merge" {
		return "", fmt.Errorf("invalid policy: %s. Allowed values are 'longest' or 'merge'", policy)
	}
	return policy, nil
}
