package main

import (
	"fmt"
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
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found or failed to load .env")
	}

	cfg := config.Load()
	fmt.Printf("Configuration loaded: %+v\n", cfg)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	stravaClient := strava.NewHttpClient(httpClient, cfg)

	save := &usecase.SaveWorkout{
		WorkoutRepo:     csv.NewCSVWorkoutRepository("workouts.csv"),
		WorkoutProvider: strava.NewProvider(stravaClient),
	}

	orches := &orchestration.SaveWorkoutsOrchestrator{
		SaveWorkoutUseCase: save,
	}

	startDate := orchestration.NewDate(2026, time.January, 5)
	endDate := orchestration.NewDate(2026, time.January, 30)
	period := orchestration.NewPeriod(startDate, endDate)

	err := orches.SaveWorkoutsOverPeriod(period, 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Workout(s) processed and saved (if any).")
	}
}
