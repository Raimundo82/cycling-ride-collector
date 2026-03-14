package auth_service

import (
	"errors"
	"testing"
	"time"

	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGoogleGetValidAccessTokenReturnsExistingCachedTokenWhenStillValid(t *testing.T) {
	Convey("Given a valid cached Google token", t, func() {
		cachedToken := &auth_model.Token{AccessToken: "cached_access_token", RefreshToken: "cached_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)
		tokenService.cachedToken = cachedToken

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the cached token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, cachedToken.AccessToken)
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsNewTokenWhenRepositoryRefreshTokenExistsAndRefreshSucceeds(t *testing.T) {
	Convey("Given repository returns a Google refresh token and provider returns a new token", t, func() {
		storedToken := &auth_model.Token{RefreshToken: "google_refresh_token"}
		newToken := &auth_model.Token{AccessToken: "new_google_access_token", RefreshToken: "new_google_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Tokens: &auth_model.Tokens{GoogleToken: storedToken}}
		tokenProvider := &mockTokenProvider{Tokens: &auth_model.Tokens{StravaToken: newToken}}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the refreshed token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, newToken.AccessToken)
				So(tokenService.cachedToken, ShouldResemble, newToken)
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRepositoryGetFails(t *testing.T) {
	Convey("Given repository returns an error when getting tokens", t, func() {
		repoErr := errors.New("repository error")

		tokenRepo := &mockTokenRepository{GetErr: repoErr}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			Convey("Then it should panic before returning the repository error", func() {
				So(func() {
					_, _ = tokenService.GetValidToken()
				}, ShouldPanic)
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRepositoryReturnsNilTokens(t *testing.T) {
	Convey("Given repository returns nil tokens", t, func() {
		tokenRepo := &mockTokenRepository{Tokens: nil}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no google tokens available")
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRepositoryReturnsNilGoogleToken(t *testing.T) {
	Convey("Given repository returns tokens without a Google token", t, func() {
		tokenRepo := &mockTokenRepository{Tokens: &auth_model.Tokens{}}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no google tokens available")
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRepositoryReturnsNoRefreshToken(t *testing.T) {
	Convey("Given repository returns a Google token without a refresh token", t, func() {
		token := &auth_model.Token{AccessToken: "stale_google_access_token", RefreshToken: "", ExpiresAt: time.Now().Add(-1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Tokens: &auth_model.Tokens{GoogleToken: token}}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the expected error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no google refresh token available")
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRefreshFails(t *testing.T) {
	Convey("Given repository returns a Google refresh token and provider returns an error", t, func() {
		refreshErr := errors.New("refresh error")
		storedToken := &auth_model.Token{RefreshToken: "google_refresh_token"}

		tokenRepo := &mockTokenRepository{Tokens: &auth_model.Tokens{GoogleToken: storedToken}}
		tokenProvider := &mockTokenProvider{Tokens: &auth_model.Tokens{}, Err: refreshErr}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the refresh error", func() {
				So(err, ShouldEqual, refreshErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGoogleGetValidAccessTokenReturnsErrorWhenRepositorySaveFailsAfterRefresh(t *testing.T) {
	Convey("Given repository returns a Google refresh token, provider returns a new token and repository fails to save", t, func() {
		saveErr := errors.New("save error")
		storedToken := &auth_model.Token{RefreshToken: "google_refresh_token"}
		newToken := &auth_model.Token{AccessToken: "new_google_access_token", RefreshToken: "new_google_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Tokens: &auth_model.Tokens{GoogleToken: storedToken}, SaveErr: saveErr}
		tokenProvider := &mockTokenProvider{Tokens: &auth_model.Tokens{StravaToken: newToken}}
		tokenService := NewGoogleTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the save error", func() {
				So(err, ShouldEqual, saveErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}
