package config

import (
	"os"
	"strconv"
)

type Config struct {
	MininmalWorkoutDuration int
}

func Load() *Config {
	val, _ := strconv.Atoi(getEnv("MINIMAL_WORKOUT_DURATION", "30"))
	return &Config{
		MininmalWorkoutDuration: val,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
