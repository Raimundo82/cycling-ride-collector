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
	stravaClientIDKey           = "STRAVA_CLIENT_ID"
	stravaClientSecretKey       = "STRAVA_CLIENT_SECRET"
	stravaRefreshTokenKey       = "STRAVA_REFRESH_TOKEN"
	stravaOAuthBaseURLKey       = "strava.oauthBaseUrl"
	googleClientIDKey           = "GOOGLE_CLIENT_ID"
	googleClientSecretKey       = "GOOGLE_CLIENT_SECRET"
	googleRefreshTokenKey       = "GOOGLE_REFRESH_TOKEN"
	googleOAuthBaseURLKey       = "googleOAuth.oauthBaseUrl"
	emailFromKey                = "EMAIL_FROM"
	emailToKey                  = "EMAIL_TO"
	emailSubjectKey             = "EMAIL_SUBJECT"
)

type Config struct {
	OutputFilePath string               `json:"outputFilePath"`
	Strava         *StravaConfig        `json:"strava"`
	GoogleOAuth    *GoogleOAuthConfig   `json:"googleOAuth"`
	Email          *EmailConfig         `json:"email"`
	ExcelTemplate  *ExcelTemplateConfig `json:"excelTemplate"`
}

type EmailConfig struct {
	From    string
	To      string
	Subject string
}

type ExcelTemplateConfig struct {
	TemplatePath string `json:"templatePath"`
	SheetName    string `json:"sheetName"`
	StartCell    string `json:"startCell"`
}

type StravaConfig struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	ApiBaseUrl   string `json:"apiBaseUrl"`
	OAuthBaseUrl string `json:"oauthBaseUrl"`
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	OAuthBaseUrl string `json:"oauthBaseUrl"`
}

func Load() (*Config, error) {
	cfg, err := loadFromFile(configFilePath)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
	}
	applySensitiveEnv(cfg)
	return cfg, err
}

func (c *Config) ValidateRequired() error {
	if c == nil {
		return fmt.Errorf("%s%s", missingRequiredValuesPrefix, strings.Join(allRequiredConfigKeys(), ", "))
	}

	missing := make([]string, 0, 11)

	missing = validateStravaConfig(c, missing)
	missing = validateGoogleOAuthConfig(c, missing)
	missing = validateEmailConfig(c, missing)

	if len(missing) > 0 {
		return fmt.Errorf("%s%s", missingRequiredValuesPrefix, strings.Join(missing, ", "))
	}

	return nil
}

func validateGoogleOAuthConfig(c *Config, missing []string) []string {
	if c.GoogleOAuth == nil || c.GoogleOAuth.ClientID == "" {
		missing = append(missing, googleClientIDKey)
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.ClientSecret == "" {
		missing = append(missing, googleClientSecretKey)
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.RefreshToken == "" {
		missing = append(missing, googleRefreshTokenKey)
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.OAuthBaseUrl == "" {
		missing = append(missing, googleOAuthBaseURLKey)
	}
	return missing
}

func validateEmailConfig(c *Config, missing []string) []string {
	if c.Email == nil || c.Email.From == "" {
		missing = append(missing, emailFromKey)
	}
	if c.Email == nil || c.Email.To == "" {
		missing = append(missing, emailToKey)
	}
	if c.Email == nil || c.Email.Subject == "" {
		missing = append(missing, emailSubjectKey)
	}
	return missing
}

func validateStravaConfig(c *Config, missing []string) []string {
	if c.Strava == nil || c.Strava.ClientId == "" {
		missing = append(missing, stravaClientIDKey)
	}
	if c.Strava == nil || c.Strava.ClientSecret == "" {
		missing = append(missing, stravaClientSecretKey)
	}
	if c.Strava == nil || c.Strava.RefreshToken == "" {
		missing = append(missing, stravaRefreshTokenKey)
	}
	if c.Strava == nil || c.Strava.OAuthBaseUrl == "" {
		missing = append(missing, stravaOAuthBaseURLKey)
	}
	return missing
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
	cfg.Strava.ClientId = getEnv("STRAVA_CLIENT_ID")
	cfg.Strava.ClientSecret = getEnv("STRAVA_CLIENT_SECRET")
	cfg.Strava.RefreshToken = getEnv("STRAVA_REFRESH_TOKEN")
	cfg.GoogleOAuth.ClientID = getEnv("GOOGLE_CLIENT_ID")
	cfg.GoogleOAuth.ClientSecret = getEnv("GOOGLE_CLIENT_SECRET")
	cfg.GoogleOAuth.RefreshToken = getEnv("GOOGLE_REFRESH_TOKEN")
	cfg.Email.From = getEnv("EMAIL_FROM")
	cfg.Email.To = getEnv("EMAIL_TO")
	cfg.Email.Subject = getEnv("EMAIL_SUBJECT")
}

func allRequiredConfigKeys() []string {
	return []string{
		stravaClientIDKey,
		stravaClientSecretKey,
		stravaRefreshTokenKey,
		stravaOAuthBaseURLKey,
		googleClientIDKey,
		googleClientSecretKey,
		googleRefreshTokenKey,
		googleOAuthBaseURLKey,
		emailFromKey,
		emailToKey,
		emailSubjectKey,
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func loadFromFile(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return initializeNestedConfigs(cfg), err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return initializeNestedConfigs(cfg), err
	}

	return initializeNestedConfigs(cfg), nil
}
