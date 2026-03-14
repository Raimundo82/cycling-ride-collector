package auth_repository

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

const getTokenFile = "test_tokens.json"

func TestGetTokensSuccess(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		stravaToken := &auth_model.Token{
			AccessToken:  "strava_valid_access_token",
			RefreshToken: "strava_valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		googleToken := &auth_model.Token{
			AccessToken:  "google_valid_access_token",
			RefreshToken: "google_valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: stravaToken, GoogleToken: googleToken})

		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			tokens, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(tokens.StravaToken.AccessToken, ShouldEqual, "strava_valid_access_token")
				So(tokens.StravaToken.RefreshToken, ShouldEqual, "strava_valid_refresh_token")
				So(tokens.StravaToken.IsExpired(), ShouldBeFalse)
				So(tokens.StravaToken.ExpiresAt, ShouldHappenAfter, time.Now())
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "google_valid_access_token")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "google_valid_refresh_token")
				So(tokens.GoogleToken.IsExpired(), ShouldBeFalse)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenAfter, time.Now())
			})
		})
	})
}

func TestGetTokensExpiredToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		stravaToken := &auth_model.Token{
			AccessToken:  "strava_access_token",
			RefreshToken: "strava_refresh_token",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		googleToken := &auth_model.Token{
			AccessToken:  "google_access_token",
			RefreshToken: "google_refresh_token",
			ExpiresAt:    time.Now().Add(-2 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: stravaToken, GoogleToken: googleToken})

		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			tokens, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(tokens.StravaToken.IsExpired(), ShouldBeTrue)
				So(tokens.StravaToken.ExpiresAt, ShouldHappenBefore, time.Now())
				So(tokens.GoogleToken.IsExpired(), ShouldBeTrue)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenBefore, time.Now())
			})
		})
	})
}

func TestGetTokensRefreshOnlyToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		stravaToken := &auth_model.Token{
			AccessToken:  "",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Time{},
		}
		googleToken := &auth_model.Token{
			AccessToken:  "",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Time{},
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: stravaToken, GoogleToken: googleToken})
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			tokens, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(tokens.StravaToken.AccessToken, ShouldEqual, "")
				So(tokens.StravaToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.StravaToken.IsExpired(), ShouldBeTrue)
				So(tokens.StravaToken.ExpiresAt, ShouldHappenBefore, time.Now())
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.GoogleToken.IsExpired(), ShouldBeTrue)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenBefore, time.Now())
			})
		})
	})
}

func TestGetGoogleTokenReturnsEmptyTokenWhenSectionMissing(t *testing.T) {
	Convey("Given a token file without a google section", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		stravaToken := &auth_model.Token{
			AccessToken:  "",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Time{},
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: stravaToken})
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)

		Convey("When getting the tokens", func() {
			tokens, err := repo.GetTokens()

			Convey("Then it returns an empty google token without error", func() {
				So(err, ShouldBeNil)
				So(tokens.StravaToken, ShouldNotBeNil)
				So(tokens.GoogleToken, ShouldNotBeNil)
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "")
			})
		})
	})
}

func TestGetTokensFileNotFound(t *testing.T) {
	Convey("Given an invalid token file and a token repo", t, func() {
		repo, err := NewTokenRepository("non_existent_file.json")
		Convey("When creating the token repo", func() {
			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(repo, ShouldBeNil)
			})
		})
	})
}

func TestGetTokensOpenFileError(t *testing.T) {
	Convey("Given a token repo pointing to a non-existent file", t, func() {
		repo := &tokenRepository{filePath: "non_existent_file.json"}

		Convey("When getting tokens", func() {
			token, err := repo.GetTokens()

			Convey("Then it returns an open file error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to open token file")
				So(token, ShouldBeNil)
			})
		})
	})
}
