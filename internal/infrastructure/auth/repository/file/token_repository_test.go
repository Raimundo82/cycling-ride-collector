package auth_repository

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	whenCreatingNewTokenRepository = "When creating new token repository"
	tokenRepositoryTestFile        = "tokens_test.json"
)

func TestNewTokenRepositorySuccess(t *testing.T) {
	Convey("Given an existing valid token file path", t, func() {
		filePath := tokenRepositoryTestFile
		tokenFile, _ := os.CreateTemp("", filePath)
		_, err := tokenFile.Write(
			[]byte(
				`{
					"strava_token": {
						"access_token":"",
						"refresh_token":"strava_refresh_token",
						"expires_at":"0001-01-01T00:00:00Z"
					},
				    "google_token":{
						"access_token":"",
						"refresh_token":"google_refresh_token",
						"expires_at":"0001-01-01T00:00:00Z"
					}
				}`),
		)
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}

		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())
			Convey("Then it should return a valid token repository", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(repo.tokens.StravaToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.RefreshToken, ShouldEqual, "strava_refresh_token")
				So(repo.tokens.StravaToken.ExpiresAt.String(), ShouldEqual, "0001-01-01 00:00:00 +0000 UTC")
				So(repo.tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.RefreshToken, ShouldEqual, "google_refresh_token")
				So(repo.tokens.GoogleToken.ExpiresAt.String(), ShouldEqual, "0001-01-01 00:00:00 +0000 UTC")
			})
		})
	})
}

func TestNewTokenRepositoryInvalidJSON(t *testing.T) {
	Convey("Given an existing invalid token file path", t, func() {
		filePath := "tokens_test.json"
		tokenFile, _ := os.CreateTemp("", filePath)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(repo, ShouldBeNil)
			})
		})
	})
}

func TestNewTokenRepositoryFileNotFound(t *testing.T) {
	Convey("Given a non-existing token file path", t, func() {
		filePath := "non_existing_tokens_test.json"

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(filePath)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to open token file")
				So(repo, ShouldBeNil)
			})
		})
	})
}

func TestEmptyTokenJsonFile(t *testing.T) {
	Convey("Given a empty token file path", t, func() {
		filePath := tokenRepositoryTestFile
		tokenFile, _ := os.CreateTemp("", filePath)
		_, err := tokenFile.Write(
			[]byte(`{}`))
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}
		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())

			Convey("Then it should return empty tokens property", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(repo.tokens, ShouldNotBeNil)
				So(repo.tokens.StravaToken, ShouldNotBeNil)
				So(repo.tokens.GoogleToken, ShouldNotBeNil)
				So(repo.tokens.StravaToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.RefreshToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.ExpiresAt.String(), ShouldEqual, "0001-01-01 00:00:00 +0000 UTC")
				So(repo.tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.RefreshToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.ExpiresAt.String(), ShouldEqual, "0001-01-01 00:00:00 +0000 UTC")
			})
		})
	})
}

func TestNewTokenRepositoryWithNullTokens(t *testing.T) {
	Convey("Given a token file with both tokens set to null", t, func() {
		tokenFile, _ := os.CreateTemp("", tokenRepositoryTestFile)
		_, err := tokenFile.Write([]byte(`{"strava_token": null, "google_token": null}`))
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}

		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())

			Convey("Then it should initialize empty token values", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(repo.tokens, ShouldNotBeNil)
				So(repo.tokens.StravaToken, ShouldNotBeNil)
				So(repo.tokens.GoogleToken, ShouldNotBeNil)
				So(repo.tokens.StravaToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.RefreshToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.RefreshToken, ShouldEqual, "")
			})
		})
	})
}

func TestNewTokenRepositoryWithMissingStravaToken(t *testing.T) {
	Convey("Given a token file with only google token defined", t, func() {
		tokenFile, _ := os.CreateTemp("", tokenRepositoryTestFile)
		_, err := tokenFile.Write([]byte(`{
			"google_token": {
				"access_token": "google_access_token",
				"refresh_token": "google_refresh_token",
				"expires_at": "0001-01-01T00:00:00Z"
			}
		}`))
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}

		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())

			Convey("Then it should initialize an empty strava token", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(repo.tokens.StravaToken, ShouldNotBeNil)
				So(repo.tokens.StravaToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.RefreshToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.AccessToken, ShouldEqual, "google_access_token")
				So(repo.tokens.GoogleToken.RefreshToken, ShouldEqual, "google_refresh_token")
			})
		})
	})
}

func TestNewTokenRepositoryWithMissingGoogleToken(t *testing.T) {
	Convey("Given a token file with only strava token defined", t, func() {
		tokenFile, _ := os.CreateTemp("", tokenRepositoryTestFile)
		_, err := tokenFile.Write([]byte(`{
			"strava_token": {
				"access_token": "strava_access_token",
				"refresh_token": "strava_refresh_token",
				"expires_at": "0001-01-01T00:00:00Z"
			}
		}`))
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}

		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey(whenCreatingNewTokenRepository, func() {
			repo, err := NewTokenRepository(tokenFile.Name())

			Convey("Then it should initialize an empty google token", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(repo.tokens.GoogleToken, ShouldNotBeNil)
				So(repo.tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(repo.tokens.GoogleToken.RefreshToken, ShouldEqual, "")
				So(repo.tokens.StravaToken.AccessToken, ShouldEqual, "strava_access_token")
				So(repo.tokens.StravaToken.RefreshToken, ShouldEqual, "strava_refresh_token")
			})
		})
	})
}
