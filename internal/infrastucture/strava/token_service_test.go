package strava

import (
	"errors"
	"testing"
	"time"

	tokens "github.com/raimundo82/cycling-ride-collector/internal/infrastucture/token_store"
	. "github.com/smartystreets/goconvey/convey"
)

type mockTokenProvider struct {
	Token *tokens.Token
	Err   error
}

// RefreshAccessToken implements [OAuthProvider].
func (m *mockTokenProvider) RefreshAccessToken(refreshToken string) (*tokens.Token, error) {
	return m.Token, m.Err
}

type mockTokenRepository struct {
	SaveErr error
	GetErr  error
	Token   *tokens.Token
}

// GetTokens implements [tokens.TokenRepository].
func (m *mockTokenRepository) GetTokens() (*tokens.Token, error) {
	return m.Token, m.GetErr
}

// SaveTokens implements [tokens.TokenRepository].
func (m *mockTokenRepository) SaveTokens(token *tokens.Token) error {
	return m.SaveErr
}

var (
	_ OAuthProvider          = (*mockTokenProvider)(nil)
	_ tokens.TokenRepository = (*mockTokenRepository)(nil)
)

func TestGetValidAccessToken_GivenValidToken_WhenExecutes_ThenReturnsToken(t *testing.T) {
	Convey("Given a valid token", t, func() {
		token := &tokens.Token{AccessToken: "access_token", RefreshToken: "refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: token}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidAccessToken()

			Convey("Then it should return the valid token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, token.AccessToken)
			})
		})
	})
}

func TestGetValidAccessToken_GivenExpiredTokenInRepoAndProviderReturnsToken_WhenExecutes_ThenReturnsNewToken(t *testing.T) {
	Convey("Given an expired token in repository and provider returns a new token", t, func() {
		expiredToken := &tokens.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &tokens.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: expiredToken}
		tokenProvider := &mockTokenProvider{Token: newToken}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidAccessToken()

			Convey("Then it should return the new token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, newToken.AccessToken)
			})
		})
	})
}

func TestGetValidAccessToken_GivenRepoReturnsError_WhenExecutes_ThenReturnsError(t *testing.T) {
	Convey("Given repository returns an error when getting token", t, func() {
		repoErr := errors.New("repository error")

		tokenRepo := &mockTokenRepository{GetErr: repoErr}
		tokenProvider := &mockTokenProvider{}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidAccessToken()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldEqual, repoErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}

func TestGetValidAccessToken_GivenExpiredTokenInRepoAndProviderReturnsTokenAndRepoSaveReturnsError_WhenExecutes_ThenReturnsError(t *testing.T) {
	Convey("Given an expired token in repository, provider returns a new token and repository returns an error when saving", t, func() {
		expiredToken := &tokens.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &tokens.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}
		saveErr := errors.New("save error")

		tokenRepo := &mockTokenRepository{Token: expiredToken, SaveErr: saveErr}
		tokenProvider := &mockTokenProvider{Token: newToken}
		tokenService := NewTokenService(tokenProvider, tokenRepo)

		Convey("When executing GetValidAccessToken", func() {
			result, err := tokenService.GetValidAccessToken()

			Convey("Then it should return the save error", func() {
				So(err, ShouldEqual, saveErr)
				So(result, ShouldBeEmpty)
			})
		})
	})
}
