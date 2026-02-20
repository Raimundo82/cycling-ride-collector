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

		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)

		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			token, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(token.AccessToken, ShouldEqual, "valid_access_token")
				So(token.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(token.IsExpired(), ShouldBeFalse)
				So(token.ExpiresAt, ShouldHappenAfter, time.Now())
			})
		})
	})
}

func TestGetTokensExpiredToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &auth_model.Token{
			AccessToken:  "invalid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			token, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(token.AccessToken, ShouldEqual, "invalid_access_token")
				So(token.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(token.IsExpired(), ShouldBeTrue)
				So(token.ExpiresAt, ShouldHappenBefore, time.Now())
			})
		})
	})
}

func TestGetTokensRefreshOnlyToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", getTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &auth_model.Token{
			AccessToken:  "",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Time{},
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When getting tokens", func() {
			token, err := repo.GetTokens()
			Convey("Then it returns the expected token", func() {
				So(err, ShouldBeNil)
				So(token.AccessToken, ShouldEqual, "")
				So(token.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(token.IsExpired(), ShouldBeTrue)
				So(token.ExpiresAt, ShouldHappenBefore, time.Now())
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
