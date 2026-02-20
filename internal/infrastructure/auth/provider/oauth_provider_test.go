package auth_provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

type mockTokenRefresher struct {
	tokenResponse *auth_model.RefreshAccessTokenResponse
	err           error
	gotRequest    *auth_model.RefreshAccessTokenRequest
}

func (m *mockTokenRefresher) RefreshToken(ctx context.Context, request *auth_model.RefreshAccessTokenRequest) (*auth_model.RefreshAccessTokenResponse, error) {
	m.gotRequest = request
	return m.tokenResponse, m.err
}

var _ TokenRefresher = (*mockTokenRefresher)(nil)

func TestOAuthProviderRefreshTokenSuccess(t *testing.T) {
	Convey("Given an oauth provider and a successful token refresher", t, func() {
		refresher := &mockTokenRefresher{
			tokenResponse: &auth_model.RefreshAccessTokenResponse{
				TokenType:    "Bearer",
				AccessToken:  "new-access-token",
				ExpiresAt:    1771095847,
				ExpiresIn:    21600,
				RefreshToken: "new-refresh-token",
			},
		}
		cfg := &config.Config{
			StravaClientId:     "test-client-id",
			StravaClientSecret: "test-client-secret",
		}
		oauthProvider := NewOAuthProvider(refresher, cfg)

		Convey("When RefreshToken is called", func() {
			token, err := oauthProvider.RefreshToken("incoming-refresh-token")

			Convey("Then it returns the mapped token and no error", func() {
				So(err, ShouldBeNil)
				So(token, ShouldNotBeNil)
				So(token.AccessToken, ShouldEqual, "new-access-token")
				So(token.RefreshToken, ShouldEqual, "new-refresh-token")
				So(token.ExpiresAt, ShouldEqual, time.Unix(1771095847, 0))
			})
		})
	})
}

func TestOAuthProviderRefreshTokenPassesExpectedRequest(t *testing.T) {
	Convey("Given an oauth provider and a token refresher", t, func() {
		refresher := &mockTokenRefresher{
			tokenResponse: &auth_model.RefreshAccessTokenResponse{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				ExpiresAt:    1771095847,
			},
		}
		cfg := &config.Config{
			StravaClientId:     "test-client-id",
			StravaClientSecret: "test-client-secret",
		}
		oauthProvider := NewOAuthProvider(refresher, cfg)

		Convey("When RefreshToken is called", func() {
			_, _ = oauthProvider.RefreshToken("incoming-refresh-token")

			Convey("Then it builds the expected request", func() {
				So(refresher.gotRequest, ShouldNotBeNil)
				So(refresher.gotRequest.ClientID, ShouldEqual, "test-client-id")
				So(refresher.gotRequest.ClientSecret, ShouldEqual, "test-client-secret")
				So(refresher.gotRequest.RefreshToken, ShouldEqual, "incoming-refresh-token")
				So(refresher.gotRequest.GrantType, ShouldEqual, "refresh_token")
			})
		})
	})
}

func TestOAuthProviderRefreshTokenReturnsRefresherError(t *testing.T) {
	Convey("Given an oauth provider and a failing token refresher", t, func() {
		refresher := &mockTokenRefresher{
			err: errors.New("strava request failed"),
		}
		cfg := &config.Config{
			StravaClientId:     "test-client-id",
			StravaClientSecret: "test-client-secret",
		}
		oauthProvider := NewOAuthProvider(refresher, cfg)

		Convey("When RefreshToken is called", func() {
			token, err := oauthProvider.RefreshToken("incoming-refresh-token")

			Convey("Then it returns the refresher error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava request failed")
				So(token, ShouldBeNil)
			})
		})
	})
}
