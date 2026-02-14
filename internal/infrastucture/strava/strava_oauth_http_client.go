package strava

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
)

type stravaOAuthHttpClient struct {
	httpClient *http.Client
	baseUrl    string
}

var _ StravaOAuthClient = (*stravaOAuthHttpClient)(nil)

func NewStravaOauthHttpClient(httpClient *http.Client, cfg *config.Config) *stravaOAuthHttpClient {
	return &stravaOAuthHttpClient{
		httpClient: httpClient,
		baseUrl:    cfg.StravaOauthBaseUrl,
	}
}

// RefreshAccessToken implements [StravaOAuthClient].
func (s *stravaOAuthHttpClient) RefreshAccessToken(request *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.baseUrl+"/token", bytes.NewBuffer(body))
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

	var refreshResponse RefreshAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return nil, err
	}

	return &refreshResponse, nil
}
