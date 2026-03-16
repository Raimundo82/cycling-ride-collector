package config

import (
	"os"
)

type Config struct {
	StravaApiBaseUrl   string
	StravaOauthBaseUrl string
	StravaClientId     string
	StravaClientSecret string
	StravaRefreshToken string
	OutputFilePath     string
	TokenFilePath      string
	Strava             *StravaConfig
	ExcelTemplate      ExcelTemplateConfig
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

func Load() *Config {
	return &Config{
		StravaApiBaseUrl:   getEnv("STRAVA_API_BASE_URL"),
		StravaOauthBaseUrl: getEnv("STRAVA_OAUTH_BASE_URL"),
		StravaClientId:     getEnv("STRAVA_CLIENT_ID"),
		StravaClientSecret: getEnv("STRAVA_CLIENT_SECRET"),
		StravaRefreshToken: getEnv("STRAVA_REFRESH_TOKEN"),
		TokenFilePath:      getEnv("TOKEN_FILE_PATH"),
		ExcelTemplate: ExcelTemplateConfig{
			TemplatePath: getEnv("EXCEL_TEMPLATE_PATH"),
			SheetName:    getEnv("EXCEL_TEMPLATE_SHEETNAME"),
			StartCell:    getEnv("EXCEL_TEMPLATE_STARTCELL"),
		},
		Strava: &StravaConfig{
			ClientId:     getEnv("STRAVA_CLIENT_ID"),
			ClientSecret: getEnv("STRAVA_CLIENT_SECRET"),
			RefreshToken: getEnv("STRAVA_REFRESH_TOKEN"),
			ApiBaseUrl:   getEnv("STRAVA_API_BASE_URL"),
			OAuthBaseUrl: getEnv("STRAVA_OAUTH_BASE_URL"),
		},
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}
