package auth_repository

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewTokenRepositorySuccess(t *testing.T) {
	Convey("Given an existing valid token file path", t, func() {
		filePath := "tokens_test.json"
		tokenFile, _ := os.CreateTemp("", filePath)
		_, err := tokenFile.Write([]byte(`{"access_token":"","refresh_token":"test_refresh_token","expires_at":"0000-01-01T00:00:00Z"}`))
		if err != nil {
			t.Fatalf("Failed to write to token file: %v", err)
		}
		_ = tokenFile.Close()
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey("When creating new token repository", func() {
			repo, err := NewTokenRepository(tokenFile.Name())
			Convey("Then it should return a valid token repository", func() {
				So(err, ShouldBeNil)
				So(repo, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(repo.token.AccessToken, ShouldEqual, "")
				So(repo.token.RefreshToken, ShouldEqual, "test_refresh_token")
				So(repo.token.ExpiresAt.String(), ShouldEqual, "0000-01-01 00:00:00 +0000 UTC")
			})
		})
	})
}

func TestNewTokenRepositoryInvalidJSON(t *testing.T) {
	Convey("Given an existing invalid token file path", t, func() {
		filePath := "tokens_test.json"
		tokenFile, _ := os.CreateTemp("", filePath)
		defer func() { _ = os.Remove(tokenFile.Name()) }()

		Convey("When creating new token repository", func() {
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

		Convey("When creating new token repository", func() {
			repo, err := NewTokenRepository(filePath)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to open token file")
				So(repo, ShouldBeNil)
			})
		})
	})
}
