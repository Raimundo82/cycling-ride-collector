package auth_service

import (
	"errors"
	"testing"
	"time"

	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

type mockTokenProvider struct {
	Token *auth_model.Token
	Err   error
}

// RefreshToken implements [authInterfaces.OAuthProvider].
func (m *mockTokenProvider) RefreshToken(refreshToken string) (*auth_model.Token, error) {
	return m.Token, m.Err
}

type mockTokenRepository struct {
	SaveErr error
	GetErr  error
	Token   *auth_model.Token
}

// GetTokens implements [authInterfaces.TokenRepository].
func (m *mockTokenRepository) GetTokens() (*auth_model.Token, error) {
	return m.Token, m.GetErr
}

// SaveTokens implements [authInterfaces.TokenRepository].
func (m *mockTokenRepository) SaveTokens(token *auth_model.Token) error {
	return m.SaveErr
}

var (
	_ authInterfaces.OAuthProvider   = (*mockTokenProvider)(nil)
	_ authInterfaces.TokenRepository = (*mockTokenRepository)(nil)
)

func TestGetValidAccessTokenReturnsExistingTokenWhenRepositoryTokenStillValid(t *testing.T) {
	Convey("Given a valid token", t, func() {
		token := &auth_model.Token{AccessToken: "access_token", RefreshToken: "refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: token}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the valid token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, token.AccessToken)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsNewTokenWhenRepositoryTokenExpiredAndRefreshSucceeds(t *testing.T) {
	Convey("Given an expired token in repository and a valid refresh token", t, func() {
		expiredToken := &auth_model.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &auth_model.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: expiredToken}
		tokenProvider := &mockTokenProvider{Token: newToken}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the new token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, newToken.AccessToken)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenTokenExpiredAndRefreshAccessTokenRequestFails(t *testing.T) {
	Convey("Given repository returns an expired token and provider returns an error when refreshing token", t, func() {
		expiredToken := &auth_model.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		refreshErr := errors.New("refresh error")

		tokenRepo := &mockTokenRepository{Token: expiredToken}
		tokenProvider := &mockTokenProvider{Err: refreshErr}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the refresh error", func() {
				So(err, ShouldEqual, refreshErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenNoAccessTokenAndRefreshAccessTokenRequestFails(t *testing.T) {
	Convey("Given repository returns no access token and provider returns an error when refreshing token", t, func() {
		expiredToken := &auth_model.Token{AccessToken: "", RefreshToken: "expired_refresh_token", ExpiresAt: time.Time{}}
		refreshErr := errors.New("refresh error")

		tokenRepo := &mockTokenRepository{Token: expiredToken}
		tokenProvider := &mockTokenProvider{Err: refreshErr}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the refresh error", func() {
				So(err, ShouldEqual, refreshErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenRepositoryGetFails(t *testing.T) {
	Convey("Given repository returns an error when getting token", t, func() {
		repoErr := errors.New("repository error")

		tokenRepo := &mockTokenRepository{GetErr: repoErr}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldEqual, repoErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenRepositoryReturnsNilTokens(t *testing.T) {
	Convey("Given repository returns nil tokens", t, func() {
		tokenRepo := &mockTokenRepository{Token: nil}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no tokens available")
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenRepositorySaveFailsAfterRefresh(t *testing.T) {
	Convey("Given an expired token in repository, provider returns a new token and repository returns an error when saving", t, func() {
		expiredToken := &auth_model.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &auth_model.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}
		saveErr := errors.New("save error")

		tokenRepo := &mockTokenRepository{Token: expiredToken, SaveErr: saveErr}
		tokenProvider := &mockTokenProvider{Token: newToken}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the save error", func() {
				So(err, ShouldEqual, saveErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenRepositoryReturnsNoValidTokenAndNoRefreshToken(t *testing.T) {
	Convey("Given repository returns empty access and refresh tokens", t, func() {
		token := &auth_model.Token{AccessToken: "", RefreshToken: "", ExpiresAt: time.Now().Add(1 * time.Hour)}
		tokenRepo := &mockTokenRepository{Token: token}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the expected error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no valid access token and no refresh token available")
			})
		})
	})
}

func TestGetValidAccessTokenReturnsErrorWhenRepositoryReturnsExpiredAccessTokenAndEmptyRefreshToken(t *testing.T) {
	Convey("Given repository returns expired access token and empty refresh token", t, func() {
		token := &auth_model.Token{AccessToken: "test-token", RefreshToken: "", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		tokenRepo := &mockTokenRepository{Token: token}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidToken()

			Convey("Then it should return the expected error", func() {
				So(err, ShouldNotBeNil)
				So(result, ShouldBeEmpty)
				So(err.Error(), ShouldEqual, "no valid access token and no refresh token available")
			})
		})
	})
}
