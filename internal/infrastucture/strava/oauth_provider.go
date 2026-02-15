package strava

import (
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	tokenStore "github.com/raimundo82/cycling-ride-collector/internal/infrastucture/token_store"
)

type OAuthProvider interface {
	RefreshAccessToken(refreshToken string) (*tokenStore.Token, error)
}

type StravaOAuthProvider struct {
	oauthClient StravaOAuthClient
	config      *config.Config
}

var _ OAuthProvider = (*StravaOAuthProvider)(nil)

func NewOAuthProvider(cfg *config.Config) OAuthProvider {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	oauthClient := NewStravaOAuthHttpClient(httpClient, cfg)
	return &StravaOAuthProvider{oauthClient: oauthClient, config: cfg}
}

// RefreshAccessToken implements [contracts.TokenProvider].
func (o *StravaOAuthProvider) RefreshAccessToken(refreshToken string) (*tokenStore.Token, error) {
	resp, err := o.oauthClient.RefreshAccessToken(&RefreshAccessTokenRequest{
		ClientID:     o.config.StravaClientId,
		ClientSecret: o.config.StravaClientSecret,
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return nil, err
	}
	return &tokenStore.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    time.Unix(int64(resp.ExpiresAt), 0),
	}, nil
}
