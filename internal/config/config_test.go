package config

import "testing"

func TestLoadShouldPopulateNestedConfigWhenEnvVarsAreSet(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "https://api.strava.test")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "https://oauth.strava.test")
	t.Setenv("STRAVA_CLIENT_ID", "client-id")
	t.Setenv("STRAVA_CLIENT_SECRET", "client-secret")
	t.Setenv("STRAVA_REFRESH_TOKEN", "refresh-token")
	t.Setenv("GOOGLE_CLIENT_ID", "google-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "google-client-secret")
	t.Setenv("GOOGLE_OAUTH_TOKEN_URL", "https://oauth2.googleapis.test/token")
	t.Setenv("EMAIL_FROM", "from@example.com")
	t.Setenv("EMAIL_TO", "to1@example.com,to2@example.com,to3@example.com")
	t.Setenv("EMAIL_SUBJECT", "Custom Subject")
	t.Setenv("EXCEL_TEMPLATE_PATH", "/tmp/template.xlsx")
	t.Setenv("EXCEL_TEMPLATE_SHEETNAME", "Workouts")
	t.Setenv("EXCEL_TEMPLATE_STARTCELL", "B8")

	cfg := Load()

	if cfg.Strava == nil {
		t.Fatalf("expected Strava config to be initialized")
	}
	if cfg.GoogleOAuth == nil {
		t.Fatalf("expected GoogleOAuth config to be initialized")
	}
	if cfg.Email == nil {
		t.Fatalf("expected Email config to be initialized")
	}
	if cfg.ExcelTemplate == nil {
		t.Fatalf("expected ExcelTemplate config to be initialized")
	}

	if cfg.Strava.ApiBaseUrl != "https://api.strava.test" {
		t.Fatalf("expected Strava.ApiBaseUrl to be loaded from env")
	}
	if cfg.Strava.OAuthBaseUrl != "https://oauth.strava.test" {
		t.Fatalf("expected Strava.OAuthBaseUrl to be loaded from env")
	}
	if cfg.Strava.ClientId != "client-id" {
		t.Fatalf("expected Strava.ClientId to be loaded from env")
	}
	if cfg.Strava.ClientSecret != "client-secret" {
		t.Fatalf("expected Strava.ClientSecret to be loaded from env")
	}
	if cfg.Strava.RefreshToken != "refresh-token" {
		t.Fatalf("expected Strava.RefreshToken to be loaded from env")
	}
	if cfg.GoogleOAuth.ClientID != "google-client-id" {
		t.Fatalf("expected GoogleOAuth.ClientID to be loaded from env")
	}
	if cfg.GoogleOAuth.ClientSecret != "google-client-secret" {
		t.Fatalf("expected GoogleOAuth.ClientSecret to be loaded from env")
	}
	if cfg.GoogleOAuth.OAuthBaseUrl != "https://oauth2.googleapis.test/token" {
		t.Fatalf("expected GoogleOAuth.OAuthBaseUrl to be loaded from env")
	}
	if cfg.Email.From != "from@example.com" {
		t.Fatalf("expected Email.From to be loaded from env")
	}
	if cfg.Email.To != "to1@example.com,to2@example.com,to3@example.com" {
		t.Fatalf("expected Email.To to be loaded from env")
	}
	if cfg.Email.Subject != "Custom Subject" {
		t.Fatalf("expected Email.Subject to be loaded from env")
	}
	if cfg.ExcelTemplate.TemplatePath != "/tmp/template.xlsx" {
		t.Fatalf("expected ExcelTemplate.TemplatePath to be loaded from env")
	}
	if cfg.ExcelTemplate.SheetName != "Workouts" {
		t.Fatalf("expected ExcelTemplate.SheetName to be loaded from env")
	}
	if cfg.ExcelTemplate.StartCell != "B8" {
		t.Fatalf("expected ExcelTemplate.StartCell to be loaded from env")
	}
}

func TestLoadShouldReturnEmptyValuesWhenEnvVarsAreMissing(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "")
	t.Setenv("STRAVA_CLIENT_ID", "")
	t.Setenv("STRAVA_CLIENT_SECRET", "")
	t.Setenv("STRAVA_REFRESH_TOKEN", "")
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("GOOGLE_OAUTH_TOKEN_URL", "")
	t.Setenv("EMAIL_FROM", "")
	t.Setenv("EMAIL_TO", "")
	t.Setenv("EMAIL_SUBJECT", "")
	t.Setenv("EXCEL_TEMPLATE_PATH", "")
	t.Setenv("EXCEL_TEMPLATE_SHEETNAME", "")
	t.Setenv("EXCEL_TEMPLATE_STARTCELL", "")

	cfg := Load()

	if cfg.Strava == nil {
		t.Fatalf("expected Strava config to be initialized")
	}
	if cfg.GoogleOAuth == nil {
		t.Fatalf("expected GoogleOAuth config to be initialized")
	}
	if cfg.Email == nil {
		t.Fatalf("expected Email config to be initialized")
	}
	if cfg.ExcelTemplate == nil {
		t.Fatalf("expected ExcelTemplate config to be initialized")
	}

	if cfg.Strava.ApiBaseUrl != "" {
		t.Fatalf("expected empty Strava.ApiBaseUrl when env var is missing")
	}
	if cfg.Strava.OAuthBaseUrl != "" {
		t.Fatalf("expected empty Strava.OAuthBaseUrl when env var is missing")
	}
	if cfg.Strava.ClientId != "" {
		t.Fatalf("expected empty Strava.ClientId when env var is missing")
	}
	if cfg.Strava.ClientSecret != "" {
		t.Fatalf("expected empty Strava.ClientSecret when env var is missing")
	}
	if cfg.Strava.RefreshToken != "" {
		t.Fatalf("expected empty Strava.RefreshToken when env var is missing")
	}
	if cfg.GoogleOAuth.ClientID != "" {
		t.Fatalf("expected empty GoogleOAuth.ClientID when env var is missing")
	}
	if cfg.GoogleOAuth.ClientSecret != "" {
		t.Fatalf("expected empty GoogleOAuth.ClientSecret when env var is missing")
	}
	if cfg.GoogleOAuth.OAuthBaseUrl != "" {
		t.Fatalf("expected empty GoogleOAuth.OAuthBaseUrl when env var is missing")
	}
	if cfg.Email.From != "" {
		t.Fatalf("expected empty Email.From when env var is missing")
	}
	if cfg.Email.To != "" {
		t.Fatalf("expected empty Email.To when env var is missing")
	}
	if cfg.Email.Subject != "" {
		t.Fatalf("expected empty Email.Subject when env var is missing")
	}
	if cfg.ExcelTemplate.TemplatePath != "" {
		t.Fatalf("expected empty ExcelTemplate.TemplatePath when env var is missing")
	}
	if cfg.ExcelTemplate.SheetName != "" {
		t.Fatalf("expected empty ExcelTemplate.SheetName when env var is missing")
	}
	if cfg.ExcelTemplate.StartCell != "" {
		t.Fatalf("expected empty ExcelTemplate.StartCell when env var is missing")
	}
	if cfg.OutputFilePath != "" {
		t.Fatalf("expected empty OutputFilePath when not set")
	}
}

func TestLoadShouldKeepGoogleRefreshTokenEmptyByDefault(t *testing.T) {
	cfg := Load()

	if cfg.GoogleOAuth == nil {
		t.Fatalf("expected GoogleOAuth config to be initialized")
	}
	if cfg.GoogleOAuth.RefreshToken != "" {
		t.Fatalf("expected GoogleOAuth.RefreshToken to be empty by default")
	}
}
