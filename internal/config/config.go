package config

import (
	"os"
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

func getEnv(key string) string {
	return os.Getenv(key)
}
