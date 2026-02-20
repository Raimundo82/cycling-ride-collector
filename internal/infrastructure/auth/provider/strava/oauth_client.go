package auth_strava

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	auth_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
)

type OAuthClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewOAuthHttpClient(httpClient *http.Client, cfg *config.Config) *OAuthClient {
	return &OAuthClient{
		httpClient: httpClient,
		baseUrl:    cfg.StravaOauthBaseUrl,
	}
}

var _ auth_provider.TokenRefresher = (*OAuthClient)(nil)

func (s *OAuthClient) RefreshToken(ctx context.Context, request *auth_model.RefreshAccessTokenRequest) (*auth_model.RefreshAccessTokenResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseUrl+"/token", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava error: %s", resp.Status)
	}

	var refreshResponse auth_model.RefreshAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return nil, err
	}

	return &refreshResponse, nil
}
