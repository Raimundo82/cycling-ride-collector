package auth_google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	auth_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
)

type OAuthClient struct {
	httpClient *http.Client
	baseUrl    string
}

const defaultGoogleTokenURL = "https://oauth2.googleapis.com/token"

func NewOAuthHttpClient(httpClient *http.Client, cfg *config.Config) *OAuthClient {
	baseURL := cfg.GoogleOAuth.TokenURL
	if baseURL == "" {
		baseURL = defaultGoogleTokenURL
	}

	return &OAuthClient{
		httpClient: httpClient,
		baseUrl:    baseURL,
	}
}

var _ auth_provider.TokenRefresher = (*OAuthClient)(nil)

func (s *OAuthClient) RefreshToken(ctx context.Context, request *auth_model.RefreshAccessTokenRequest) (*auth_model.RefreshAccessTokenResponse, error) {
	form := url.Values{}
	form.Set("client_id", request.ClientID)
	form.Set("client_secret", request.ClientSecret)
	form.Set("grant_type", request.GrantType)
	form.Set("refresh_token", request.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google error: %s", resp.Status)
	}

	var refreshResponse auth_model.RefreshAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return nil, err
	}

	return &refreshResponse, nil
}
