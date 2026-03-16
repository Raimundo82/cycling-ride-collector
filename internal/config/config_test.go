package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "https://api.strava.test")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "https://oauth.strava.test")
	t.Setenv("STRAVA_CLIENT_ID", "client-id")
	t.Setenv("STRAVA_CLIENT_SECRET", "client-secret")
	t.Setenv("STRAVA_REFRESH_TOKEN", "refresh-token")
	t.Setenv("TOKEN_FILE_PATH", "/tmp/token.json")
	t.Setenv("GOOGLE_CLIENT_ID", "google-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "google-client-secret")
	t.Setenv("GOOGLE_OAUTH_TOKEN_URL", "https://oauth2.googleapis.test/token")
	t.Setenv("EMAIL_FROM", "from@example.com")
	t.Setenv("EMAIL_TO", "to1@example.com,to2@example.com,to3@example.com")
	t.Setenv("EMAIL_SUBJECT", "Custom Subject")

	cfg := Load()

	if cfg.StravaApiBaseUrl != "https://api.strava.test" {
		t.Fatalf("expected StravaApiBaseUrl to be loaded from env")
	}
	if cfg.StravaOauthBaseUrl != "https://oauth.strava.test" {
		t.Fatalf("expected StravaOauthBaseUrl to be loaded from env")
	}
	if cfg.StravaClientId != "client-id" {
		t.Fatalf("expected StravaClientId to be loaded from env")
	}
	if cfg.StravaClientSecret != "client-secret" {
		t.Fatalf("expected StravaClientSecret to be loaded from env")
	}
	if cfg.StravaRefreshToken != "refresh-token" {
		t.Fatalf("expected StravaRefreshToken to be loaded from env")
	}
	if cfg.TokenFilePath != "/tmp/token.json" {
		t.Fatalf("expected TokenFilePath to be loaded from env")
	}
	if cfg.GoogleOAuth.ClientID != "google-client-id" {
		t.Fatalf("expected GoogleOAuth.ClientID to be loaded from env")
	}
	if cfg.GoogleOAuth.ClientSecret != "google-client-secret" {
		t.Fatalf("expected GoogleOAuth.ClientSecret to be loaded from env")
	}
	if cfg.GoogleOAuth.TokenURL != "https://oauth2.googleapis.test/token" {
		t.Fatalf("expected GoogleOAuth.TokenURL to be loaded from env")
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
}

func TestLoadMissingEnv(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "")
	t.Setenv("STRAVA_CLIENT_ID", "")
	t.Setenv("STRAVA_CLIENT_SECRET", "")
	t.Setenv("STRAVA_REFRESH_TOKEN", "")
	t.Setenv("TOKEN_FILE_PATH", "")
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	t.Setenv("GOOGLE_OAUTH_TOKEN_URL", "")
	t.Setenv("EMAIL_FROM", "")
	t.Setenv("EMAIL_TO", "")
	t.Setenv("EMAIL_SUBJECT", "")

	cfg := Load()

	if cfg.StravaApiBaseUrl != "" {
		t.Fatalf("expected empty StravaApiBaseUrl when env var is missing")
	}
	if cfg.StravaOauthBaseUrl != "" {
		t.Fatalf("expected empty StravaOauthBaseUrl when env var is missing")
	}
	if cfg.StravaClientId != "" {
		t.Fatalf("expected empty StravaClientId when env var is missing")
	}
	if cfg.StravaClientSecret != "" {
		t.Fatalf("expected empty StravaClientSecret when env var is missing")
	}
	if cfg.StravaRefreshToken != "" {
		t.Fatalf("expected empty StravaRefreshToken when env var is missing")
	}
	if cfg.TokenFilePath != "" {
		t.Fatalf("expected empty TokenFilePath when env var is missing")
	}
	if cfg.GoogleOAuth.ClientID != "" {
		t.Fatalf("expected empty GoogleOAuth.ClientID when env var is missing")
	}
	if cfg.GoogleOAuth.ClientSecret != "" {
		t.Fatalf("expected empty GoogleOAuth.ClientSecret when env var is missing")
	}
	if cfg.GoogleOAuth.TokenURL != "" {
		t.Fatalf("expected default GoogleOAuth.TokenURL when env var is missing")
	}
	if cfg.Email.From != "" {
		t.Fatalf("expected empty Email.From when env var is missing")
	}
	if cfg.Email.To != "" {
		t.Fatalf("expected empty Email.To when env var is missing")
	}
	if cfg.Email.Subject != "" {
		t.Fatalf("expected default Email.Subject when env var is missing")
	}
}
