package auth

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	BaseURL      string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}
	return Config{
		ClientID:     getEnv("STRAVA_CLIENT_ID", ""),
		ClientSecret: getEnv("STRAVA_CLIENT_SECRET", ""),
		RedirectURI:  getEnv("STRAVA_REDIRECT_URI", "http://localhost:8080/callback"),
		BaseURL:      getEnv("STRAVA_BASE_URL", "https://www.strava.com"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
