package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "https://api.strava.test")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "https://oauth.strava.test")
	t.Setenv("STRAVA_CLIENT_ID", "client-id")
	t.Setenv("STRAVA_CLIENT_SECRET", "client-secret")
	t.Setenv("TOKEN_FILE_PATH", "/tmp/token.json")

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
	if cfg.TokenFilePath != "/tmp/token.json" {
		t.Fatalf("expected TokenFilePath to be loaded from env")
	}
}

func TestLoadMissingEnv(t *testing.T) {
	t.Setenv("STRAVA_API_BASE_URL", "")
	t.Setenv("STRAVA_OAUTH_BASE_URL", "")
	t.Setenv("STRAVA_CLIENT_ID", "")
	t.Setenv("STRAVA_CLIENT_SECRET", "")
	t.Setenv("TOKEN_FILE_PATH", "")

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
	if cfg.TokenFilePath != "" {
		t.Fatalf("expected empty TokenFilePath when env var is missing")
	}
}
