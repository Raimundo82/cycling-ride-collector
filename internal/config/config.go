package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	OutputFilePath string
	Strava         *StravaConfig
	GoogleOAuth    *GoogleOAuthConfig
	Email          *EmailConfig
	ExcelTemplate  *ExcelTemplateConfig
}

type EmailConfig struct {
	From    string
	To      string
	Subject string
}

type ExcelTemplateConfig struct {
	TemplatePath string
	SheetName    string
	StartCell    string
}

type StravaConfig struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
	ApiBaseUrl   string
	OAuthBaseUrl string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	OAuthBaseUrl string
}

func Load() *Config {
	return &Config{
		Strava: &StravaConfig{
			ClientId:     getEnv("STRAVA_CLIENT_ID"),
			ClientSecret: getEnv("STRAVA_CLIENT_SECRET"),
			RefreshToken: getEnv("STRAVA_REFRESH_TOKEN"),
			ApiBaseUrl:   getEnv("STRAVA_API_BASE_URL"),
			OAuthBaseUrl: getEnv("STRAVA_OAUTH_BASE_URL"),
		},
		GoogleOAuth: &GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
			OAuthBaseUrl: getEnv("GOOGLE_OAUTH_TOKEN_URL"),
			RefreshToken: getEnv("GOOGLE_REFRESH_TOKEN"),
		},
		Email: &EmailConfig{
			From:    getEnv("EMAIL_FROM"),
			To:      getEnv("EMAIL_TO"),
			Subject: getEnv("EMAIL_SUBJECT"),
		},
		ExcelTemplate: &ExcelTemplateConfig{
			TemplatePath: getEnv("EXCEL_TEMPLATE_PATH"),
			SheetName:    getEnv("EXCEL_TEMPLATE_SHEETNAME"),
			StartCell:    getEnv("EXCEL_TEMPLATE_STARTCELL"),
		},
	}
}

func (c *Config) ValidateRequired() error {
	missing := make([]string, 0)

	if c == nil {
		return fmt.Errorf("missing required config values: STRAVA_API_BASE_URL, STRAVA_OAUTH_BASE_URL, STRAVA_CLIENT_ID, STRAVA_CLIENT_SECRET, STRAVA_REFRESH_TOKEN, GOOGLE_OAUTH_TOKEN_URL, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GOOGLE_REFRESH_TOKEN, EMAIL_FROM, EMAIL_TO, EMAIL_SUBJECT")
	}

	if c.Strava == nil || c.Strava.ApiBaseUrl == "" {
		missing = append(missing, "STRAVA_API_BASE_URL")
	}
	if c.Strava == nil || c.Strava.OAuthBaseUrl == "" {
		missing = append(missing, "STRAVA_OAUTH_BASE_URL")
	}
	if c.Strava == nil || c.Strava.ClientId == "" {
		missing = append(missing, "STRAVA_CLIENT_ID")
	}
	if c.Strava == nil || c.Strava.ClientSecret == "" {
		missing = append(missing, "STRAVA_CLIENT_SECRET")
	}
	if c.Strava == nil || c.Strava.RefreshToken == "" {
		missing = append(missing, "STRAVA_REFRESH_TOKEN")
	}

	if c.GoogleOAuth == nil || c.GoogleOAuth.OAuthBaseUrl == "" {
		missing = append(missing, "GOOGLE_OAUTH_TOKEN_URL")
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.ClientID == "" {
		missing = append(missing, "GOOGLE_CLIENT_ID")
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.ClientSecret == "" {
		missing = append(missing, "GOOGLE_CLIENT_SECRET")
	}
	if c.GoogleOAuth == nil || c.GoogleOAuth.RefreshToken == "" {
		missing = append(missing, "GOOGLE_REFRESH_TOKEN")
	}

	if c.Email == nil || c.Email.From == "" {
		missing = append(missing, "EMAIL_FROM")
	}
	if c.Email == nil || c.Email.To == "" {
		missing = append(missing, "EMAIL_TO")
	}
	if c.Email == nil || c.Email.Subject == "" {
		missing = append(missing, "EMAIL_SUBJECT")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config values: %s", strings.Join(missing, ", "))
	}

	return nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}
