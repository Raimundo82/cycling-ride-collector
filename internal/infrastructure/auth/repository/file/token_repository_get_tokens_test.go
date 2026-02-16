package file

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTokensSuccess(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &model.Token{
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
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &model.Token{
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
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &model.Token{
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

func TestGetTokensInvalidJSON(t *testing.T) {
	Convey("Given an invalid token file content and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		_, _ = tokenFile.WriteString(`{"access_token":"valid_access_token","refresh_token":"valid_refresh_token","expires_at":"0000-01-01T00:00:00Z"}`)
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		_ = os.WriteFile(tokenFile.Name(), []byte("invalid json content"), 0o600)
		Convey("When getting tokens", func() {
			token, err := repo.GetTokens()
			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(token, ShouldBeNil)
			})
		})
	})
}
