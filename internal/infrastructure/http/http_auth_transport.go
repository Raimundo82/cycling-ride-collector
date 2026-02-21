package custom_http

import (
	"fmt"
	"net/http"

	activity_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/interfaces"
)

type authTransport struct {
	underlying    http.RoundTripper
	tokenProvider activity_interfaces.TokenProvider
}

var _ http.RoundTripper = (*authTransport)(nil)

// RoundTrip implements [http.RoundTripper].
func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := a.tokenProvider.GetValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token for request: %w", err)
	}
	if accessToken != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return a.underlying.RoundTrip(req)
}

func NewAuthHttpClient(tokenProvider activity_interfaces.TokenProvider) *http.Client {
	return &http.Client{
		Transport: &authTransport{
			underlying:    http.DefaultTransport,
			tokenProvider: tokenProvider,
		},
	}
}
