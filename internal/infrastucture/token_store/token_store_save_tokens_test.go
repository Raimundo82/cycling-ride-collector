package tokens

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSaveTokens_GivenValidTokenFileAndTokenStoreAndValidToken_WhenSavingTokens_ThenItSavesTheExpectedToken(t *testing.T) {
	Convey("Given a valid token file and, a token and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = tokenFile.Close()

		store := &TokenStore{filePath: tokenFile.Name()}
		Convey("When saving tokens", func() {
			err := store.SaveTokens(token)

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

func TestSaveTokens_GivenInvalidTokenFileAndTokenStoreAndValidToken_WhenSavingTokens_ThenItReturnsAnError(t *testing.T) {
	Convey("Given an invalid token file and, a token and a token store", t, func() {
		store := &TokenStore{filePath: "/invalid/path/token.json"}
		token := &Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}

		Convey("When saving tokens", func() {
			err := store.SaveTokens(token)

			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSaveTokens_GivenValidTokenFileAndTokenStoreAndInvalidToken_WhenSavingTokens_ThenItSavesTheExpectedToken(t *testing.T) {
	Convey("Given a valid token file and, an invalid token and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()
		store := &TokenStore{filePath: tokenFile.Name()}
		token := &Token{
			AccessToken:  "",
			RefreshToken: "",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		Convey("When saving tokens", func() {
			err := store.SaveTokens(token)
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

func readTokenFromFile(filepath string) (*Token, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	var token Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}
