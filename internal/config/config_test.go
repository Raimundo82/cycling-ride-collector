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
		"apiBaseUrl": "https://api.strava.test",
		"oauthBaseUrl": "https://oauth.strava.test"
	},
	"googleOAuth": {
		"oauthBaseUrl": "https://oauth2.googleapis.test/token"
	},
	"email": {
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
		setRuntimeEnv(
			t,
			"strava-client-id",
			"client-secret",
			"refresh-token",
			"google-client-id",
			"google-client-secret",
			"google-refresh-token",
			"from@example.com",
			"to@example.com",
		)

		Convey("When Load is called", func() {
			cfg, err := Load()

			Convey("Then JSON values should come from the file", func() {
				So(err, ShouldBeNil)
				So(cfg.OutputFilePath, ShouldEqual, "workouts.xlsx")
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.Strava.ApiBaseUrl, ShouldEqual, "https://api.strava.test")
				So(cfg.Strava.OAuthBaseUrl, ShouldEqual, "https://oauth.strava.test")
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.GoogleOAuth.OAuthBaseUrl, ShouldEqual, "https://oauth2.googleapis.test/token")
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.Email.Subject, ShouldEqual, "Workout report")
				So(cfg.ExcelTemplate, ShouldNotBeNil)
				So(cfg.ExcelTemplate.TemplatePath, ShouldEqual, "template.xlsx")
				So(cfg.ExcelTemplate.SheetName, ShouldEqual, "Sheet1")
				So(cfg.ExcelTemplate.StartCell, ShouldEqual, "A1")
			})

			Convey("Then sensitive values should come from environment variables", func() {
				So(cfg.Strava.ClientId, ShouldEqual, "strava-client-id")
				So(cfg.Strava.ClientSecret, ShouldEqual, "client-secret")
				So(cfg.Strava.RefreshToken, ShouldEqual, "refresh-token")
				So(cfg.GoogleOAuth.ClientID, ShouldEqual, "google-client-id")
				So(cfg.GoogleOAuth.ClientSecret, ShouldEqual, "google-client-secret")
				So(cfg.GoogleOAuth.RefreshToken, ShouldEqual, "google-refresh-token")
				So(cfg.Email.From, ShouldEqual, "from@example.com")
				So(cfg.Email.To, ShouldEqual, "to@example.com")
			})
		})
	})
}

func TestLoadShouldFailValidationWhenSensitiveEnvVarsAreMissing(t *testing.T) {
	Convey("Given a config.json file and missing secret environment variables", t, func() {
		writeConfigFixture(t, fullConfigFixture)
		setRuntimeEnv(t, "", "", "", "", "", "", "", "")

		Convey("When Load is called", func() {
			cfg, err := Load()

			Convey("Then JSON values should still be loaded", func() {
				So(err, ShouldBeNil)
				So(cfg.OutputFilePath, ShouldEqual, "workouts.xlsx")
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.Strava.ApiBaseUrl, ShouldEqual, "https://api.strava.test")
				So(cfg.Strava.OAuthBaseUrl, ShouldEqual, "https://oauth.strava.test")
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.GoogleOAuth.OAuthBaseUrl, ShouldEqual, "https://oauth2.googleapis.test/token")
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.Email.Subject, ShouldEqual, "Workout report")
			})

			Convey("Then sensitive values should be empty before validation", func() {
				So(cfg.Strava.ClientId, ShouldEqual, "")
				So(cfg.Strava.ClientSecret, ShouldEqual, "")
				So(cfg.Strava.RefreshToken, ShouldEqual, "")
				So(cfg.GoogleOAuth.ClientID, ShouldEqual, "")
				So(cfg.GoogleOAuth.ClientSecret, ShouldEqual, "")
				So(cfg.GoogleOAuth.RefreshToken, ShouldEqual, "")
				So(cfg.Email.From, ShouldEqual, "")
				So(cfg.Email.To, ShouldEqual, "")
			})

			Convey("Then validation should fail because sensitive values are mandatory", func() {
				err := cfg.ValidateRequired()

				So(err, ShouldNotBeNil)
				for _, key := range []string{
					stravaClientIDKey,
					stravaClientSecretKey,
					stravaRefreshTokenKey,
					googleClientIDKey,
					googleClientSecretKey,
					googleRefreshTokenKey,
					emailFromKey,
					emailToKey,
				} {
					So(err.Error(), ShouldContainSubstring, key)
				}
			})
		})
	})
}

func TestLoadShouldInitializeNestedConfigsWhenJSONIsEmptyObject(t *testing.T) {
	Convey("Given a config.json file with an empty json object", t, func() {
		writeConfigFixture(t, `{}`)
		setRuntimeEnv(t, "", "", "", "", "", "", "", "")

		Convey("When Load is called", func() {
			cfg, err := Load()

			Convey("Then nested config objects should still be initialized", func() {
				So(err, ShouldBeNil)
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

func TestLoadShouldReturnErrorWhenConfigFileDoesNotExist(t *testing.T) {
	Convey("Given a working directory without config.json", t, func() {
		changeWorkingDir(t, t.TempDir())
		setRuntimeEnv(t, "", "", "", "", "", "", "", "")

		Convey("When Load is called", func() {
			cfg, err := Load()

			Convey("Then it should return an error and still initialize nested configs", func() {
				So(err, ShouldNotBeNil)
				So(cfg, ShouldNotBeNil)
				So(cfg.Strava, ShouldNotBeNil)
				So(cfg.GoogleOAuth, ShouldNotBeNil)
				So(cfg.Email, ShouldNotBeNil)
				So(cfg.ExcelTemplate, ShouldNotBeNil)
			})
		})
	})
}

func TestValidateRequiredShouldReturnErrorWhenConfigIsNil(t *testing.T) {
	Convey("Given a nil config", t, func() {
		var cfg *Config

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should return an error listing all required config keys", func() {
				So(err, ShouldNotBeNil)
				for _, key := range allRequiredConfigKeys() {
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

			Convey("Then it should return an error listing the missing config keys", func() {
				So(err, ShouldNotBeNil)
				for _, key := range allRequiredConfigKeys() {
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
				ClientId:     "strava-client-id",
				ClientSecret: "strava-client-secret",
				RefreshToken: "strava-refresh-token",
				OAuthBaseUrl: "https://www.strava.com/oauth/token",
				ApiBaseUrl:   "https://www.strava.com/api/v3",
			},
			GoogleOAuth: &GoogleOAuthConfig{
				ClientID:     "google-client-id",
				ClientSecret: "google-client-secret",
				RefreshToken: "google-refresh-token",
				OAuthBaseUrl: "https://oauth2.googleapis.com/token",
			},
			Email: &EmailConfig{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Workout report",
			},
			ExcelTemplate: &ExcelTemplateConfig{
				SheetName: "Sheet1",
				StartCell: "A1",
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

func setRuntimeEnv(t *testing.T, stravaID, stravaSecret, stravaRefresh, googleID, googleSecret, googleRefresh, emailFrom, emailTo string) {
	t.Helper()
	t.Setenv(stravaClientIDKey, stravaID)
	t.Setenv(stravaClientSecretKey, stravaSecret)
	t.Setenv(stravaRefreshTokenKey, stravaRefresh)
	t.Setenv(googleClientIDKey, googleID)
	t.Setenv(googleClientSecretKey, googleSecret)
	t.Setenv(googleRefreshTokenKey, googleRefresh)
	t.Setenv(emailFromKey, emailFrom)
	t.Setenv(emailToKey, emailTo)
}

func TestAllRequiredConfigKeysShouldReturnTheExpectedKeys(t *testing.T) {
	Convey("Given the required config keys list", t, func() {
		keys := allRequiredConfigKeys()

		Convey("Then it should contain the required JSON and env-backed keys", func() {
			So(keys, ShouldResemble, []string{
				stravaClientIDKey,
				stravaClientSecretKey,
				stravaRefreshTokenKey,
				stravaOAuthBaseURLKey,
				stravaAPIBaseURLKey,
				googleClientIDKey,
				googleClientSecretKey,
				googleRefreshTokenKey,
				googleOAuthBaseURLKey,
				emailFromKey,
				emailToKey,
			})
		})
	})
}

func TestValidateRequiredErrorShouldNotMentionOptionalConfigKeys(t *testing.T) {
	Convey("Given a config with missing required values", t, func() {
		cfg := &Config{}

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should not mention optional config keys", func() {
				So(err, ShouldNotBeNil)
				So(strings.Contains(err.Error(), "outputFilePath"), ShouldBeFalse)
				So(strings.Contains(err.Error(), "excelTemplate.templatePath"), ShouldBeFalse)
			})
		})
	})
}

func TestValidateRequiredShouldReturnErrorWhenOAuthFieldsAreMissing(t *testing.T) {
	Convey("Given a config with env-backed values but missing OAuth urls", t, func() {
		cfg := &Config{
			Strava: &StravaConfig{
				ClientId:     "strava-client-id",
				ClientSecret: "strava-client-secret",
				RefreshToken: "strava-refresh-token",
			},
			GoogleOAuth: &GoogleOAuthConfig{
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

		Convey("When ValidateRequired is called", func() {
			err := cfg.ValidateRequired()

			Convey("Then it should fail fast on the missing OAuth urls", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, stravaOAuthBaseURLKey)
				So(err.Error(), ShouldContainSubstring, googleOAuthBaseURLKey)
			})
		})
	})
}
