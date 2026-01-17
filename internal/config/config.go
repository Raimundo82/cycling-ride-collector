package config

import "os"

// Config holds application configuration
type Config struct {
	StravaAPIKey      string
	SheetsSpreadsheetID string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		StravaAPIKey:      getEnv("STRAVA_API_KEY", ""),
		SheetsSpreadsheetID: getEnv("SHEETS_SPREADSHEET_ID", ""),
	}
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
