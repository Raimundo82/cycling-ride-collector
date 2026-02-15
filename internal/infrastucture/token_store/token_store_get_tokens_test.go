package tokens

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTokens_GivenValidTokenFileAndTokenStoreAndValidToken_WhenGettingTokens_ThenItReturnsTheExpectedToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &Token{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)

		_ = tokenFile.Close()

		store := &TokenStore{filePath: tokenFile.Name()}
		Convey("When getting tokens", func() {
			token, err := store.GetTokens()
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

func TestGetTokens_GivenValidTokenFileAndTokenStoreAndExpiredToken_WhenGettingTokens_ThenItReturnsTheExpectedToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &Token{
			AccessToken:  "invalid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Now().Add(-1 * time.Hour),
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)
		_ = tokenFile.Close()

		store := &TokenStore{filePath: tokenFile.Name()}
		Convey("When getting tokens", func() {
			token, err := store.GetTokens()
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

func TestGetTokens_GivenValidTokenFileAndTokenStoreAndRefreshToken_WhenGettingTokens_ThenItReturnsTheExpectedToken(t *testing.T) {
	Convey("Given a valid token file and, a dto token and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		token := &Token{
			AccessToken:  "",
			RefreshToken: "valid_refresh_token",
			ExpiresAt:    time.Time{},
		}
		_ = json.NewEncoder(tokenFile).Encode(*token)
		_ = tokenFile.Close()

		store := &TokenStore{filePath: tokenFile.Name()}
		Convey("When getting tokens", func() {
			token, err := store.GetTokens()
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

func TestGetTokens_GivenInvalidTokenFileAndTokenStore_WhenGettingTokens_ThenItReturnsAnError(t *testing.T) {
	Convey("Given an invalid token file and a token store", t, func() {
		store := &TokenStore{filePath: "non_existent_file.json"}
		Convey("When getting tokens", func() {
			token, err := store.GetTokens()
			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(token, ShouldBeNil)
			})
		})
	})
}

func TestGetTokens_GivenInvalidTokenFileContentAndTokenStore_WhenGettingTokens_ThenItReturnsAnError(t *testing.T) {
	Convey("Given an invalid token file content and a token store", t, func() {
		tokenFile, _ := os.CreateTemp("", "token.json")
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		_, _ = tokenFile.WriteString("invalid json content")
		_ = tokenFile.Close()

		store := &TokenStore{filePath: tokenFile.Name()}
		Convey("When getting tokens", func() {
			token, err := store.GetTokens()
			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
				So(token, ShouldBeNil)
			})
		})
	})
}
