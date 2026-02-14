package strava

import (
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
)

var _ contracts.TokenProvider = (*OAuthProvider)(nil)

type OAuthProvider struct {
	oauthClient StravaOAuthClient
	config      *config.Config
}

func NewOAuthProvider(cfg *config.Config) *OAuthProvider {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	oauthClient := NewStravaOAuthHttpClient(httpClient, cfg)
	return &OAuthProvider{oauthClient: oauthClient, config: cfg}
}

// RefreshAccessToken implements [contracts.TokenProvider].
func (o *OAuthProvider) RefreshAccessToken(refreshToken string) (*dto.Token, error) {
	resp, err := o.oauthClient.RefreshAccessToken(&RefreshAccessTokenRequest{
		ClientID:     o.config.StravaClientId,
		ClientSecret: o.config.StravaClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return nil, err
	}
	return &dto.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    time.Unix(int64(resp.ExpiresAt), 0),
	}, nil
}
