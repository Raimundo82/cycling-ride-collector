package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
	. "github.com/smartystreets/goconvey/convey"
)

type mockTokenProvider struct {
	Token *dto.Token
	Err   error
}

type mockTokenRepository struct {
	SaveErr error
	GetErr  error
	Token   *dto.Token
}

// RefreshAccessToken implements [contracts.TokenProvider].
func (m *mockTokenProvider) RefreshAccessToken() (*dto.Token, error) {
	return m.Token, m.Err
}

// GetTokens implements [contracts.TokenRepository].
func (m *mockTokenRepository) GetTokens() (*dto.Token, error) {
	return m.Token, m.GetErr
}

// SaveTokens implements [contracts.TokenRepository].
func (m *mockTokenRepository) SaveTokens(token *dto.Token) error {
	return m.SaveErr
}

var _ contracts.TokenProvider = (*mockTokenProvider)(nil)
var _ contracts.TokenRepository = (*mockTokenRepository)(nil)

func TestGetAccessToken_GivenValidToken_WhenExecutes_ThenReturnsToken(t *testing.T) {
	Convey("Given a valid token", t, func() {
		token := &dto.Token{AccessToken: "access_token", RefreshToken: "refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: token}
		tokenProvider := &mockTokenProvider{}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the valid token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, token)
			})
		})

	})
}

func TestGetAccessToken_GivenNoTokenInRepoAndProviderReturnsToken_WhenExecutes_ThenReturnsNewToken(t *testing.T) {
	Convey("Given no token in repository and provider returns a new token", t, func() {
		newToken := &dto.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: nil}
		tokenProvider := &mockTokenProvider{Token: newToken}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the new token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, newToken)
			})
		})
	})
}

func TestGetAccessToken_GivenExpiredTokenInRepoAndProviderReturnsToken_WhenExecutes_ThenReturnsNewToken(t *testing.T) {
	Convey("Given an expired token in repository and provider returns a new token", t, func() {
		expiredToken := &dto.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &dto.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}

		tokenRepo := &mockTokenRepository{Token: expiredToken}
		tokenProvider := &mockTokenProvider{Token: newToken}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the new token without error", func() {
				So(err, ShouldBeNil)
				So(result, ShouldResemble, newToken)
			})
		})
	})
}

func TestGetAccessToken_GivenNoTokenInRepoAndProviderReturnsError_WhenExecutes_ThenReturnsError(t *testing.T) {
	Convey("Given no token in repository and provider returns an error", t, func() {
		providerErr := errors.New("provider error")

		tokenRepo := &mockTokenRepository{Token: nil}
		tokenProvider := &mockTokenProvider{Err: providerErr}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the provider error", func() {
				So(err, ShouldEqual, providerErr)
				So(result, ShouldBeNil)
			})
		})
	})
}

func TestGetAccessToken_GivenRepoReturnsError_WhenExecutes_ThenReturnsError(t *testing.T) {
	Convey("Given repository returns an error when getting token", t, func() {
		repoErr := errors.New("repository error")

		tokenRepo := &mockTokenRepository{GetErr: repoErr}
		tokenProvider := &mockTokenProvider{}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the repository error", func() {
				So(err, ShouldEqual, repoErr)
				So(result, ShouldBeNil)
			})
		})
	})
}

func TestGetAccessToken_GivenExpiredTokenInRepoAndProviderReturnsTokenAndRepoSaveReturnsError_WhenExecutes_ThenReturnsError(t *testing.T) {
	Convey("Given an expired token in repository, provider returns a new token and repository returns an error when saving", t, func() {
		expiredToken := &dto.Token{AccessToken: "expired_access_token", RefreshToken: "expired_refresh_token", ExpiresAt: time.Now().Add(-1 * time.Hour)}
		newToken := &dto.Token{AccessToken: "new_access_token", RefreshToken: "new_refresh_token", ExpiresAt: time.Now().Add(1 * time.Hour)}
		saveErr := errors.New("save error")

		tokenRepo := &mockTokenRepository{Token: expiredToken, SaveErr: saveErr}
		tokenProvider := &mockTokenProvider{Token: newToken}
		getAccessToken := &GetAccessToken{TokenProvider: tokenProvider, TokenRepository: tokenRepo}

		Convey("When executing GetAccessToken", func() {
			result, err := getAccessToken.Execute()

			Convey("Then it should return the save error", func() {
				So(err, ShouldEqual, saveErr)
				So(result, ShouldBeNil)
			})
		})
	})
}
