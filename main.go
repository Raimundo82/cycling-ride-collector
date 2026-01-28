package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
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
	err := save.Execute(time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC), 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Workout(s) processed and saved (if any).")
	}
}
