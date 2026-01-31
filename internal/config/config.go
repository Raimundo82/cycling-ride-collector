package config

import (
	"os"
	"strconv"
)

type Config struct {
	MinimalWorkoutDuration int
	StravaApiBaseUrl       string
	StravaAccessToken      string
}

func Load() *Config {
	const defaultDuration = 30

	valStr := getEnv("MINIMAL_WORKOUT_DURATION", "30")
	val, err := strconv.Atoi(valStr)
	if err != nil || val <= 0 {
		val = defaultDuration
	}

	return &Config{
		MinimalWorkoutDuration: val,
		StravaApiBaseUrl:       getEnv("STRAVA_API_BASE_URL", ""),
		StravaAccessToken:      getEnv("STRAVA_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
