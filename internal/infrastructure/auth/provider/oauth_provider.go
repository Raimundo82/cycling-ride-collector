package provider

import (
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type oauthProvider struct {
	oauthClient TokenRefresher
	config      *config.Config
}

var _ interfaces.OAuthProvider = (*oauthProvider)(nil)

// RefreshToken implements [interfaces.OAuthProvider].
func (o *oauthProvider) RefreshToken(refreshToken string) (*model.Token, error) {
	resp, err := o.oauthClient.RefreshToken(&RefreshAccessTokenRequest{
		ClientID:     o.config.StravaClientId,
		ClientSecret: o.config.StravaClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return nil, err
	}
	return &model.Token{
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
