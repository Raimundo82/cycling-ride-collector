package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/config"
)

type httpClient struct {
	httpClient  *http.Client
	baseUrl     string
	accessToken string
	cfg         *config.Config
}

var _ client = (*httpClient)(nil)

func NewHttpClient(http *http.Client, cfg *config.Config) *httpClient {
	return &httpClient{
		httpClient:  http,
		baseUrl:     cfg.StravaApiBaseUrl,
		accessToken: cfg.StravaAccessToken,
		cfg:         cfg,
	}
}

func (c *httpClient) addAuthHeader(req *http.Request) {
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
}

func (c *httpClient) GetActivitiesByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(time.Hour * 24)
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+"/athlete/activities", nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeader(req)

	q := req.URL.Query()
	q.Set("after", fmt.Sprint(start.Unix()))
	q.Set("before", fmt.Sprint(end.Unix()))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava error: %s", resp.Status)
	}

	var acts []*ActivityDto
	if err := json.NewDecoder(resp.Body).Decode(&acts); err != nil {
		return nil, err
	}

	return acts, nil
}

func (c *httpClient) GetWattsStream(ctx context.Context, id int64) (*WattsStreamDto, error) {
	u := fmt.Sprintf("%s/activities/%d/streams?keys=watts&key_by_type=true", c.baseUrl, id)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava error: %s", resp.Status)
	}

	var streams wattsStreamResponse
	if err := json.NewDecoder(resp.Body).Decode(&streams); err != nil {
		return &WattsStreamDto{WattsData: []int{}}, nil
	}

	if len(streams.Watts.WattsData) == 0 {
		return &WattsStreamDto{WattsData: []int{}}, nil
	}

	return &streams.Watts, nil
}

// RefreshTokenResponse represents the response from the Strava token refresh endpoint
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshAccessToken refreshes the access token using the refresh token
func (c *httpClient) RefreshAccessToken(ctx context.Context) (*RefreshTokenResponse, error) {
	tokenURL := c.cfg.StravaBaseUrl + "/oauth/token"

	// Create form data
	formData := url.Values{}
	formData.Set("client_id", c.cfg.StravaClientID)
	formData.Set("client_secret", c.cfg.StravaClientSecret)
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", c.cfg.StravaRefreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava token refresh error: %s", resp.Status)
	}

	var tokenResp RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Update the client's access token
	c.accessToken = tokenResp.AccessToken

	return &tokenResp, nil
}
