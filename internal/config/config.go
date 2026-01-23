package config

import (
	"os"
	"strconv"
)

type Config struct {
	MinimalWorkoutDuration int
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
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
