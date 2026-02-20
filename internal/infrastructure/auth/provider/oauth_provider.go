package auth_provider

import (
	"context"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type oauthProvider struct {
	oauthClient TokenRefresher
	config      *config.Config
}

var _ auth_interfaces.OAuthProvider = (*oauthProvider)(nil)

// RefreshToken implements [auth_interfaces.OAuthProvider].
func (o *oauthProvider) RefreshToken(refreshToken string) (*auth_model.Token, error) {
	resp, err := o.oauthClient.RefreshToken(
		context.Background(),
		&auth_model.RefreshAccessTokenRequest{
			ClientID:     o.config.StravaClientId,
			ClientSecret: o.config.StravaClientSecret,
			RefreshToken: refreshToken,
			GrantType:    "refresh_token",
		})
	if err != nil {
		return nil, err
	}
	return &auth_model.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    time.Unix(int64(resp.ExpiresAt), 0),
	}, nil
}

func NewOAuthProvider(oauthClient TokenRefresher, cfg *config.Config) *oauthProvider {
	return &oauthProvider{
		oauthClient: oauthClient,
		config:      cfg,
	}
}
