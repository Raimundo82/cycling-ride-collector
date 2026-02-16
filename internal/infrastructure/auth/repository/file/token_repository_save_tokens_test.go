package file

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSaveTokensSuccess(t *testing.T) {
	Convey("Given a valid token file and, a token and a token repo", t, func() {
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
		Convey("When saving tokens", func() {
			err := repo.SaveTokens(token)

			Convey("Then it saves the expected token", func() {
				So(err, ShouldBeNil)
				savedToken, err := readTokenFromFile(tokenFile.Name())
				So(err, ShouldBeNil)
				So(savedToken.AccessToken, ShouldEqual, "valid_access_token")
				So(savedToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(savedToken.IsExpired(), ShouldBeFalse)
				So(savedToken.ExpiresAt, ShouldHappenAfter, time.Now())
			})
		})
	})
}

func TestSaveTokensWriteFileError(t *testing.T) {
	Convey("Given an invalid token file and, a token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		token := &model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)
		_ = tokenFile.Close()
		_ = os.Chmod(tokenFile.Name(), 0o400)
		defer func() { _ = os.Chmod(tokenFile.Name(), 0o600) }()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)

		Convey("When saving tokens", func() {
			err := repo.SaveTokens(token)

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to write token file")
			})
		})
	})
}

func TestSaveTokensExpiredToken(t *testing.T) {
	Convey("Given a valid token file and, an invalid token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		initialToken := &model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*initialToken)
		_ = tokenFile.Close()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		token := &model.Token{
			AccessToken:  "",
			RefreshToken: "",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		Convey("When saving tokens", func() {
			err := repo.SaveTokens(token)
			Convey("Then it saves the expected token", func() {
				So(err, ShouldBeNil)
				savedToken, err := readTokenFromFile(tokenFile.Name())
				So(err, ShouldBeNil)
				So(savedToken.AccessToken, ShouldEqual, "")
				So(savedToken.RefreshToken, ShouldEqual, "")
				So(savedToken.IsExpired(), ShouldBeTrue)
				So(savedToken.ExpiresAt, ShouldHappenBefore, time.Now())
			})
		})
	})
}

func TestSaveTokensSerializationError(t *testing.T) {
	Convey("Given a token with an unmarshalable time and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		initialToken := &model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*initialToken)
		_ = tokenFile.Close()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		token := &model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		Convey("When saving tokens", func() {
			err := repo.SaveTokens(token)

			Convey("Then it returns a serialization error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to serialize token")
				So(err.Error(), ShouldContainSubstring, "year outside of range [0,9999]")
			})
		})
	})
}

func TestSaveTokensSetFilePermissionsError(t *testing.T) {
	Convey("Given a token repo writing to a device file", t, func() {
		repo := &tokenRepository{filePath: "/dev/null"}
		token := &model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		Convey("When saving tokens", func() {
			err := repo.SaveTokens(token)

			Convey("Then it returns a file permission error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to set token file permissions")
			})
		})
	})
}

func readTokenFromFile(filepath string) (*model.Token, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	var token model.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
