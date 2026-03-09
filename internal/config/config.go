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
	ExcelTemplate      ExcelTemplateConfig
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
