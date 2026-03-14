package config

import (
	"os"
)

type Config struct {
	StravaApiBaseUrl   string
	StravaOauthBaseUrl string
	StravaClientId     string
	StravaClientSecret string
	OutputFilePath     string
	TokenFilePath      string
	GoogleOAuth        GoogleOAuthConfig
	Email              EmailConfig
	ExcelTemplate      ExcelTemplateConfig
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
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

func Load() *Config {
	return &Config{
		StravaApiBaseUrl:   getEnv("STRAVA_API_BASE_URL"),
		StravaOauthBaseUrl: getEnv("STRAVA_OAUTH_BASE_URL"),
		StravaClientId:     getEnv("STRAVA_CLIENT_ID"),
		StravaClientSecret: getEnv("STRAVA_CLIENT_SECRET"),
		TokenFilePath:      getEnv("TOKEN_FILE_PATH"),
		GoogleOAuth: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
			TokenURL:     getEnv("GOOGLE_OAUTH_TOKEN_URL"),
		},
		Email: EmailConfig{
			From:    getEnv("EMAIL_FROM"),
			To:      getEnv("EMAIL_TO"),
			Subject: getEnv("EMAIL_SUBJECT"),
		},
		ExcelTemplate: ExcelTemplateConfig{
			TemplatePath: getEnv("EXCEL_TEMPLATE_PATH"),
			SheetName:    getEnv("EXCEL_TEMPLATE_SHEETNAME"),
			StartCell:    getEnv("EXCEL_TEMPLATE_STARTCELL"),
		},
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
