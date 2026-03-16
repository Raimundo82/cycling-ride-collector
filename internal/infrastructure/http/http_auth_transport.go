package custom_http

import (
	"fmt"
	"net/http"

	auth_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
)

type authTransport struct {
	underlying    http.RoundTripper
	tokenProvider auth_provider.TokenProvider
}

var _ http.RoundTripper = (*authTransport)(nil)

// RoundTrip implements [http.RoundTripper].
func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if a.tokenProvider != nil {
		accessToken, err := a.tokenProvider.GetValidToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get valid token for request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)

	}
	return a.underlying.RoundTrip(req)
}

func NewAuthTransport(tokenProvider auth_provider.TokenProvider) *authTransport {
	return &authTransport{
		underlying:    http.DefaultTransport,
		tokenProvider: tokenProvider,
	}
}
