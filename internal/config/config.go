package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	MinimalWorkoutDuration int
	StravaApiBaseUrl       string
	StravaAccessToken      string
	StartDate              time.Time
	EndDate                time.Time
	OutputFilePath         string
	DailyWorkoutPolicy     string
}

func Load() *Config {
	minimalWorkoutDurationValue, err := strconv.Atoi(getEnv("MINIMAL_WORKOUT_DURATION"))
	if err != nil || minimalWorkoutDurationValue <= 0 {
		minimalWorkoutDurationValue = 30
	}

	return &Config{
		MinimalWorkoutDuration: minimalWorkoutDurationValue,
		StravaApiBaseUrl:       getEnv("STRAVA_API_BASE_URL"),
		StravaAccessToken:      getEnv("STRAVA_ACCESS_TOKEN"),
		OutputFilePath:         getEnv("OUTPUT_FILE_PATH"),
		DailyWorkoutPolicy:     getEnv("DAILY_WORKOUT_POLICY"),
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
