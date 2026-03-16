package config

import (
	"strings"
	"testing"
)

func TestLoadShouldPopulateNestedConfigWhenEnvVarsAreSet(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "https://api.strava.test")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "https://oauth.strava.test")
	t.Setenv("STRAVA_CLIENT_ID", "client-id")
	t.Setenv("STRAVA_CLIENT_SECRET", "client-secret")
	t.Setenv("STRAVA_REFRESH_TOKEN", "refresh-token")
	t.Setenv("GOOGLE_CLIENT_ID", "google-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "google-client-secret")
	t.Setenv("GOOGLE_OAUTH_TOKEN_URL", "https://oauth2.googleapis.test/token")
	t.Setenv("GOOGLE_REFRESH_TOKEN", "google-refresh-token")
	t.Setenv("EMAIL_FROM", "from@example.com")
	t.Setenv("EMAIL_TO", "to@example.com")
	t.Setenv("EMAIL_SUBJECT", "Workout report")

	cfg := Load()

	if cfg.Strava == nil {
		t.Fatalf("expected Strava config to be initialized")
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
	if cfg.GoogleOAuth == nil {
		t.Fatalf("expected GoogleOAuth config to be initialized")
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
	if cfg.GoogleOAuth.RefreshToken != "google-refresh-token" {
		t.Fatalf("expected GoogleOAuth.RefreshToken to be loaded from env")
	}
	if cfg.Email == nil {
		t.Fatalf("expected Email config to be initialized")
	}
	if cfg.Email.From != "from@example.com" {
		t.Fatalf("expected Email.From to be loaded from env")
	}
	if cfg.Email.To != "to@example.com" {
		t.Fatalf("expected Email.To to be loaded from env")
	}
	if cfg.Email.Subject != "Workout report" {
		t.Fatalf("expected Email.Subject to be loaded from env")
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
	t.Setenv("GOOGLE_REFRESH_TOKEN", "")
	t.Setenv("EMAIL_FROM", "")
	t.Setenv("EMAIL_TO", "")
	t.Setenv("EMAIL_SUBJECT", "")

	cfg := Load()

	if cfg.Strava == nil {
		t.Fatalf("expected Strava config to be initialized")
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
	if cfg.GoogleOAuth == nil {
		t.Fatalf("expected GoogleOAuth config to be initialized")
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
	if cfg.GoogleOAuth.RefreshToken != "" {
		t.Fatalf("expected empty GoogleOAuth.RefreshToken when env var is missing")
	}
	if cfg.Email == nil {
		t.Fatalf("expected Email config to be initialized")
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
}

func TestValidateRequiredReturnsErrorWithConfigIsNil(t *testing.T) {
	var cfg *Config

	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatalf("expected validation error when config is nil")
	}
}

func TestValidateRequiredReturnsErrorWithMissingValues(t *testing.T) {
	cfg := &Config{}

	err := cfg.ValidateRequired()
	if err == nil {
		t.Fatalf("expected validation error when required config values are missing")
	}

	for _, expected := range []string{
		"STRAVA_API_BASE_URL",
		"STRAVA_OAUTH_BASE_URL",
		"STRAVA_CLIENT_ID",
		"STRAVA_CLIENT_SECRET",
		"STRAVA_REFRESH_TOKEN",
		"GOOGLE_OAUTH_TOKEN_URL",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		"GOOGLE_REFRESH_TOKEN",
		"EMAIL_FROM",
		"EMAIL_TO",
		"EMAIL_SUBJECT",
	} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("expected validation error to contain %q, got %q", expected, err.Error())
		}
	}
}

func TestValidateRequiredReturnsNilWhenAllValuesExist(t *testing.T) {
	cfg := &Config{
		Strava: &StravaConfig{
			ApiBaseUrl:   "https://www.strava.com/api/v3",
			OAuthBaseUrl: "https://www.strava.com/oauth/token",
			ClientId:     "strava-client-id",
			ClientSecret: "strava-client-secret",
			RefreshToken: "strava-refresh-token",
		},
		GoogleOAuth: &GoogleOAuthConfig{
			OAuthBaseUrl: "https://oauth2.googleapis.com/token",
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
		},
		Email: &EmailConfig{
			From:    "from@example.com",
			To:      "to@example.com",
			Subject: "Workout report",
		},
	}

	err := cfg.ValidateRequired()
	if err != nil {
		t.Fatalf("expected no validation error, got %v", err)
	}
}
