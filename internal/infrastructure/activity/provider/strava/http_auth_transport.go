package strava

import (
	"net/http"

	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/interfaces"
)

type authTransport struct {
	underlying    http.RoundTripper
	tokenProvider interfaces.TokenProvider
}

var _ http.RoundTripper = (*authTransport)(nil)

// RoundTrip implements [http.RoundTripper].
func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := a.tokenProvider.GetValidToken()
	if err != nil {
		return nil, err
	}
	if accessToken != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	return a.underlying.RoundTrip(req)
}
