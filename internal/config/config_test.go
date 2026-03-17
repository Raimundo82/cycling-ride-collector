package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const fullConfigFixture = `{
	"outputFilePath": "workouts.xlsx",
	"strava": {
		"clientId": "client-id",
		"apiBaseUrl": "https://api.strava.test",
		"oauthBaseUrl": "https://oauth.strava.test"
	},
	"googleOAuth": {
		"clientId": "google-client-id",
		"oauthBaseUrl": "https://oauth2.googleapis.test/token"
	},
	"email": {
		"from": "from@example.com",
		"to": "to@example.com",
		"subject": "Workout report"
	},
	"excelTemplate": {
		"templatePath": "template.xlsx",
		"sheetName": "Sheet1",
		"startCell": "A1"
	}
}`

func TestLoadShouldPopulateConfigFromJSONAndSensitiveFieldsFromEnv(t *testing.T) {
	Convey("Given a config.json file and secrets in environment variables", t, func() {
		writeConfigFixture(t, fullConfigFixture)
		setSensitiveEnv(t, "client-secret", "refresh-token", "google-client-secret", "google-refresh-token")

		Convey("When Load is called", func() {
			cfg := Load()

			Convey("Then JSON values should come from the file", func() {
				So(cfg.OutputFilePath, ShouldEqual, "workouts.xlsx")
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.Strava.ClientId, ShouldEqual, "client-id")
				So(cfg.Strava.ApiBaseUrl, ShouldEqual, "https://api.strava.test")
				So(cfg.Strava.OAuthBaseUrl, ShouldEqual, "https://oauth.strava.test")
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.GoogleOAuth.ClientID, ShouldEqual, "google-client-id")
				So(cfg.GoogleOAuth.OAuthBaseUrl, ShouldEqual, "https://oauth2.googleapis.test/token")
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.Email.From, ShouldEqual, "from@example.com")
				So(cfg.Email.To, ShouldEqual, "to@example.com")
				So(cfg.Email.Subject, ShouldEqual, "Workout report")
				So(cfg.ExcelTemplate, ShouldNotBeNil)
				So(cfg.ExcelTemplate.TemplatePath, ShouldEqual, "template.xlsx")
				So(cfg.ExcelTemplate.SheetName, ShouldEqual, "Sheet1")
				So(cfg.ExcelTemplate.StartCell, ShouldEqual, "A1")
			})

			Convey("Then sensitive values should come from environment variables", func() {
				So(cfg.Strava.ClientSecret, ShouldEqual, "client-secret")
				So(cfg.Strava.RefreshToken, ShouldEqual, "refresh-token")
				So(cfg.GoogleOAuth.ClientSecret, ShouldEqual, "google-client-secret")
				So(cfg.GoogleOAuth.RefreshToken, ShouldEqual, "google-refresh-token")
			})
		})
	})
}

func TestLoadShouldFailValidationWhenSensitiveEnvVarsAreMissing(t *testing.T) {
	Convey("Given a config.json file and missing secret environment variables", t, func() {
		writeConfigFixture(t, fullConfigFixture)
		setSensitiveEnv(t, "", "", "", "")

		Convey("When Load is called", func() {
			cfg := Load()

			Convey("Then JSON values should still be loaded", func() {
				So(cfg.OutputFilePath, ShouldEqual, "workouts.xlsx")
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.Strava.ClientId, ShouldEqual, "client-id")
				So(cfg.Strava.ApiBaseUrl, ShouldEqual, "https://api.strava.test")
				So(cfg.Strava.OAuthBaseUrl, ShouldEqual, "https://oauth.strava.test")
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.GoogleOAuth.ClientID, ShouldEqual, "google-client-id")
				So(cfg.GoogleOAuth.OAuthBaseUrl, ShouldEqual, "https://oauth2.googleapis.test/token")
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.Email.From, ShouldEqual, "from@example.com")
				So(cfg.Email.To, ShouldEqual, "to@example.com")
				So(cfg.Email.Subject, ShouldEqual, "Workout report")
			})

			Convey("Then sensitive values should be empty before validation", func() {
				So(cfg.Strava.ClientSecret, ShouldEqual, "")
				So(cfg.Strava.RefreshToken, ShouldEqual, "")
				So(cfg.GoogleOAuth.ClientSecret, ShouldEqual, "")
				So(cfg.GoogleOAuth.RefreshToken, ShouldEqual, "")
			})

			Convey("Then validation should fail because sensitive values are mandatory", func() {
				err := cfg.ValidateRequired()

				So(err, ShouldNotBeNil)
				for _, key := range allRequiredEnvKeys() {
					So(err.Error(), ShouldContainSubstring, key)
				}
			})
		})
	})
}

func TestLoadShouldInitializeNestedConfigsWhenJSONIsEmptyObject(t *testing.T) {
	Convey("Given a config.json file with an empty json object", t, func() {
		writeConfigFixture(t, `{}`)
		setSensitiveEnv(t, "", "", "", "")

		Convey("When Load is called", func() {
			cfg := Load()

			Convey("Then nested config objects should still be initialized", func() {
				So(cfg, ShouldNotBeNil)
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.ExcelTemplate, ShouldNotBeNil)
			})

			Convey("Then initialized nested objects should contain zero values", func() {
				So(cfg.OutputFilePath, ShouldEqual, "")
				So(cfg.Strava.ClientId, ShouldEqual, "")
				So(cfg.Strava.ApiBaseUrl, ShouldEqual, "")
				So(cfg.Strava.OAuthBaseUrl, ShouldEqual, "")
				So(cfg.GoogleOAuth.ClientID, ShouldEqual, "")
				So(cfg.GoogleOAuth.OAuthBaseUrl, ShouldEqual, "")
				So(cfg.Email.From, ShouldEqual, "")
				So(cfg.Email.To, ShouldEqual, "")
				So(cfg.Email.Subject, ShouldEqual, "")
				So(cfg.ExcelTemplate.TemplatePath, ShouldEqual, "")
				So(cfg.ExcelTemplate.SheetName, ShouldEqual, "")
				So(cfg.ExcelTemplate.StartCell, ShouldEqual, "")
			})
		})
	})
}

func TestValidateRequiredShouldReturnErrorWhenConfigIsNil(t *testing.T) {
	Convey("Given a nil config", t, func() {
		var cfg *Config

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should return an error listing all required env keys", func() {
				So(err, ShouldNotBeNil)
				for _, key := range allRequiredEnvKeys() {
					So(err.Error(), ShouldContainSubstring, key)
				}
			})
		})
	})
}

func TestValidateRequiredShouldReturnErrorWhenSensitiveValuesAreMissing(t *testing.T) {
	Convey("Given a config without the required sensitive values", t, func() {
		cfg := &Config{}

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should return an error listing the missing env keys", func() {
				So(err, ShouldNotBeNil)
				for _, key := range allRequiredEnvKeys() {
					So(err.Error(), ShouldContainSubstring, key)
				}
			})
		})
	})
}

func TestValidateRequiredShouldReturnNilWhenSensitiveValuesExist(t *testing.T) {
	Convey("Given a config with all required sensitive values", t, func() {
		cfg := &Config{
			Strava: &StravaConfig{
				ClientSecret: "strava-client-secret",
				RefreshToken: "strava-refresh-token",
			},
			GoogleOAuth: &GoogleOAuthConfig{
				ClientSecret: "google-client-secret",
				RefreshToken: "google-refresh-token",
			},
		}

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should return no error", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func writeConfigFixture(t *testing.T, content string) {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, configFilePath)
	err := os.WriteFile(configPath, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("expected config fixture to be written, got %v", err)
	}

	changeWorkingDir(t, tempDir)
}

func changeWorkingDir(t *testing.T, dir string) {
	t.Helper()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("expected working directory, got %v", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("expected to chdir to fixture dir, got %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(currentDir)
	})
}

func setSensitiveEnv(t *testing.T, stravaSecret, stravaRefresh, googleSecret, googleRefresh string) {
	t.Helper()
	t.Setenv(stravaClientSecretKey, stravaSecret)
	t.Setenv(stravaRefreshTokenKey, stravaRefresh)
	t.Setenv(googleClientSecretKey, googleSecret)
	t.Setenv(googleRefreshTokenKey, googleRefresh)
}

func TestAllRequiredEnvKeysShouldReturnTheExpectedKeys(t *testing.T) {
	Convey("Given the required env keys list", t, func() {
		keys := allRequiredEnvKeys()

		Convey("Then it should contain the four sensitive env keys", func() {
			So(keys, ShouldResemble, []string{
				stravaClientSecretKey,
				stravaRefreshTokenKey,
				googleClientSecretKey,
				googleRefreshTokenKey,
			})
		})
	})
}

func TestValidateRequiredErrorShouldOnlyMentionConfiguredEnvKeys(t *testing.T) {
	Convey("Given a config with missing sensitive values", t, func() {
		cfg := &Config{}

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should not mention JSON-backed fields", func() {
				So(err, ShouldNotBeNil)
				So(strings.Contains(err.Error(), "STRAVA_API_BASE_URL"), ShouldBeFalse)
				So(strings.Contains(err.Error(), "GOOGLE_CLIENT_ID"), ShouldBeFalse)
				So(strings.Contains(err.Error(), "EMAIL_FROM"), ShouldBeFalse)
			})
		})
	})
}
