package config

import (
	"os"
	"strconv"
)

type Config struct {
	MinimalWorkoutDuration int
	StravaBaseUrl          string
	StravaApiBaseUrl       string
	StravaClientID         string
	StravaClientSecret     string
	StravaAccessToken      string
	StravaRefreshToken     string
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
		StravaBaseUrl:          getEnv("STRAVA_BASE_URL", ""),
		StravaApiBaseUrl:       getEnv("STRAVA_API_BASE_URL", ""),
		StravaClientID:         getEnv("STRAVA_CLIENT_ID", ""),
		StravaClientSecret:     getEnv("STRAVA_CLIENT_SECRET", ""),
		StravaAccessToken:      getEnv("STRAVA_ACCESS_TOKEN", ""),
		StravaRefreshToken:     getEnv("STRAVA_REFRESH_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
