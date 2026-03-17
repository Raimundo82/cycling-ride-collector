package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	configFilePath              = "config.json"
	missingRequiredValuesPrefix = "missing required config values: "
	stravaClientSecretKey       = "STRAVA_CLIENT_SECRET"
	stravaRefreshTokenKey       = "STRAVA_REFRESH_TOKEN"
	googleClientSecretKey       = "GOOGLE_CLIENT_SECRET"
	googleRefreshTokenKey       = "GOOGLE_REFRESH_TOKEN"
)

type Config struct {
	OutputFilePath string               `json:"outputFilePath"`
	Strava         *StravaConfig        `json:"strava"`
	GoogleOAuth    *GoogleOAuthConfig   `json:"googleOAuth"`
	Email          *EmailConfig         `json:"email"`
	ExcelTemplate  *ExcelTemplateConfig `json:"excelTemplate"`
}

type EmailConfig struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}

type ExcelTemplateConfig struct {
	TemplatePath string `json:"templatePath"`
	SheetName    string `json:"sheetName"`
	StartCell    string `json:"startCell"`
}

type StravaConfig struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken"`
	ApiBaseUrl   string `json:"apiBaseUrl"`
	OAuthBaseUrl string `json:"oauthBaseUrl"`
}

type GoogleOAuthConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken"`
	OAuthBaseUrl string `json:"oauthBaseUrl"`
}

func Load() *Config {
	cfg := loadFromFile(configFilePath)
	applySensitiveEnv(cfg)
	return cfg
}

func (c *Config) ValidateRequired() error {
	if c == nil {
		return fmt.Errorf("%s%s", missingRequiredValuesPrefix, strings.Join(allRequiredEnvKeys(), ", "))
	}

	missing := make([]string, 0, 4)

	if c.Strava == nil || c.Strava.ClientSecret == "" {
		missing = append(missing, stravaClientSecretKey)
	}
	if c.Strava == nil || c.Strava.RefreshToken == "" {
		missing = append(missing, stravaRefreshTokenKey)
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.ClientSecret == "" {
		missing = append(missing, googleClientSecretKey)
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.RefreshToken == "" {
		missing = append(missing, googleRefreshTokenKey)
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s%s", missingRequiredValuesPrefix, strings.Join(missing, ", "))
	}

	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func loadFromFile(path string) *Config {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err == nil && len(data) > 0 {
		_ = json.Unmarshal(data, cfg)
	}

	return initializeNestedConfigs(cfg)
}

func initializeNestedConfigs(cfg *Config) *Config {
	if cfg.Strava == nil {
		cfg.Strava = &StravaConfig{}
	}
	if cfg.GoogleOAuth == nil {
		cfg.GoogleOAuth = &GoogleOAuthConfig{}
	}
	if cfg.Email == nil {
		cfg.Email = &EmailConfig{}
	}
	if cfg.ExcelTemplate == nil {
		cfg.ExcelTemplate = &ExcelTemplateConfig{}
	}

	return cfg
}

func applySensitiveEnv(cfg *Config) {
	cfg.Strava.ClientSecret = getEnv("STRAVA_CLIENT_SECRET")
	cfg.Strava.RefreshToken = getEnv("STRAVA_REFRESH_TOKEN")
	cfg.GoogleOAuth.ClientSecret = getEnv("GOOGLE_CLIENT_SECRET")
	cfg.GoogleOAuth.RefreshToken = getEnv("GOOGLE_REFRESH_TOKEN")
}

func allRequiredEnvKeys() []string {
	return []string{
		stravaClientSecretKey,
		stravaRefreshTokenKey,
		googleClientSecretKey,
		googleRefreshTokenKey,
	}
}
