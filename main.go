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
	// Load .env file
	_ = godotenv.Load()
	cfg := config.Load()

	startDateStr := flag.String("start-date", "", "Start date in MM/DD/YYYY format")
	endDateStr := flag.String("end-date", "", "End date in MM/DD/YYYY format")
	accessToken := flag.String("access-token", "", "Strava API access token")
	minimalWorkoutDuration := flag.Int("min-duration", cfg.MinimalWorkoutDuration, "Minimal workout duration in minutes")
	flag.Parse()

	if *accessToken != "" {
		cfg.StravaAccessToken = *accessToken
	}
	if flag.Lookup("min-duration").Value.String() != fmt.Sprint(cfg.MinimalWorkoutDuration) {
		cfg.MinimalWorkoutDuration = *minimalWorkoutDuration
	}

	// fmt.Printf("Configuration loaded: %+v\n", cfg)

	// Parse dates (from flags only, or you can add to config if needed)
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

	fmt.Printf("Start: %s\nEnd: %s\nAccess Token: %s\nMin Workout Duration: %d\n",
		startDate.Format(layout), endDate.Format(layout), cfg.StravaAccessToken, cfg.MinimalWorkoutDuration)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	stravaClient := strava.NewHttpClient(httpClient, cfg)

	save := &usecase.SaveWorkout{
		WorkoutRepo:     csv.NewCSVWorkoutRepository("workouts.csv"),
		WorkoutProvider: strava.NewProvider(stravaClient),
	}

	orches := &orchestration.SaveWorkoutsOrchestrator{
		SaveWorkoutUseCase: save,
	}
	period := orchestration.NewPeriod(startDate, endDate)

	orchestrationError := orches.SaveWorkoutsOverPeriod(period, cfg.MinimalWorkoutDuration)
	if orchestrationError != nil {
		fmt.Printf("Error: %v\n", orchestrationError)
	} else {
		fmt.Println("Workout(s) processed and saved (if any).")
	}
}
