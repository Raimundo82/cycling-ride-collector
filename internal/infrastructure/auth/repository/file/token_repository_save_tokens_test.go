package auth_repository

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

const saveTokenFile = "test_tokens.json"

func TestSaveStravaTokenSuccessAndKeepGoogleToken(t *testing.T) {
	Convey("Given a valid token file and, a strava token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", saveTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{GoogleToken: token})
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When saving strava token", func() {
			err := repo.SaveStravaToken(token)

			Convey("Then it saves the expected token", func() {
				So(err, ShouldBeNil)
				tokens, err := readTokenFromFile(tokenFile.Name())
				So(err, ShouldBeNil)
				So(tokens.StravaToken.AccessToken, ShouldEqual, "valid_access_token")
				So(tokens.StravaToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.StravaToken.IsExpired(), ShouldBeFalse)
				So(tokens.StravaToken.ExpiresAt, ShouldHappenAfter, time.Now())
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "valid_access_token")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.GoogleToken.IsExpired(), ShouldBeFalse)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenAfter, time.Now())
			})
		})
	})
}

func TestSaveGoogleTokenSuccessAndKeepStravaToken(t *testing.T) {
	Convey("Given a valid token file and, a google token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", saveTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: token})
		_ = tokenFile.Close()

		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		Convey("When saving google token", func() {
			err := repo.SaveGoogleToken(token)

			Convey("Then it saves the expected token", func() {
				So(err, ShouldBeNil)
				tokens, err := readTokenFromFile(tokenFile.Name())
				So(err, ShouldBeNil)
				So(tokens.StravaToken.AccessToken, ShouldEqual, "valid_access_token")
				So(tokens.StravaToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.StravaToken.IsExpired(), ShouldBeFalse)
				So(tokens.StravaToken.ExpiresAt, ShouldHappenAfter, time.Now())
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "valid_access_token")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "valid_refresh_token")
				So(tokens.GoogleToken.IsExpired(), ShouldBeFalse)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenAfter, time.Now())
			})
		})
	})
}

func TestSaveGoogleTokenWriteFileError(t *testing.T) {
	Convey("Given an invalid token file and, a token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", saveTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{GoogleToken: token})
		_ = tokenFile.Close()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		repo.filePath = tokenFile.Name() + "/missing/tokens.json"

		Convey("When saving google token", func() {
			err := repo.SaveGoogleToken(token)

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to write token file")
			})
		})
	})
}

func TestSaveTokensExpiredToken(t *testing.T) {
	Convey("Given a valid token file and, an invalid token and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", saveTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		initialToken := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{GoogleToken: initialToken})
		_ = tokenFile.Close()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		token := &auth_model.Token{
			AccessToken:  "",
			RefreshToken: "",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		Convey("When saving  a google token", func() {
			err := repo.SaveGoogleToken(token)
			Convey("Then it saves the expected token", func() {
				So(err, ShouldBeNil)
				tokens, err := readTokenFromFile(tokenFile.Name())
				So(err, ShouldBeNil)
				So(tokens.GoogleToken.AccessToken, ShouldEqual, "")
				So(tokens.GoogleToken.RefreshToken, ShouldEqual, "")
				So(tokens.GoogleToken.IsExpired(), ShouldBeTrue)
				So(tokens.GoogleToken.ExpiresAt, ShouldHappenBefore, time.Now())
			})
		})
	})
}

func TestSaveTokensSerializationError(t *testing.T) {
	Convey("Given a token with an unmarshalable time and a token repo", t, func() {
		tokenFile, _ := os.CreateTemp("", saveTokenFile)
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		initialToken := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(&auth_model.Tokens{StravaToken: initialToken})
		_ = tokenFile.Close()
		repo, err := NewTokenRepository(tokenFile.Name())
		So(err, ShouldBeNil)
		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		Convey("When saving tokens", func() {
			err := repo.SaveGoogleToken(token)

			Convey("Then it returns a serialization error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to serialize tokens")
				So(err.Error(), ShouldContainSubstring, "year outside of range [0,9999]")
			})
		})
	})
}

func TestSaveTokensSetFilePermissionsError(t *testing.T) {
	Convey("Given a token repo writing to a device file", t, func() {
		repo := &tokenRepository{
			filePath: "/dev/null",
			tokens:   ensureTokens(nil),
		}
		token := &auth_model.Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		Convey("When saving tokens", func() {
			err := repo.SaveStravaToken(token)

			Convey("Then it returns a file permission error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to write token file")
			})
		})
	})
}

func readTokenFromFile(filepath string) (*auth_model.Tokens, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	var tokens auth_model.Tokens
	if err := json.NewDecoder(file).Decode(&tokens); err != nil {
		return nil, err
	}
	return &tokens, nil
}
