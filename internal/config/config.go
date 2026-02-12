package config

import (
	"os"
)

type Config struct {
	StravaApiBaseUrl  string
	StravaAccessToken string
	OutputFilePath    string
}

func Load() *Config {
	return &Config{
		StravaApiBaseUrl:  getEnv("STRAVA_API_BASE_URL"),
		StravaAccessToken: getEnv("STRAVA_ACCESS_TOKEN"),
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
