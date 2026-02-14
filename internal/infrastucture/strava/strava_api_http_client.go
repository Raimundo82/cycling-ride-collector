package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
)

type stravaApiHttpClient struct {
	httpClient  *http.Client
	baseUrl     string
	accessToken string
}

type authTransport struct {
	underlying  http.RoundTripper
	accessToken string
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if a.accessToken != "" {
		req = req.Clone(req.Context())
		req.Header.Set("Authorization", "Bearer "+a.accessToken)
	}
	return a.underlying.RoundTrip(req)
}

var (
	_ StravaApiClient   = (*stravaApiHttpClient)(nil)
	_ http.RoundTripper = (*authTransport)(nil)
)

func NewStravaApiHttpClient(httpClient *http.Client, cfg *config.Config) *stravaApiHttpClient {
	if httpClient.Transport == nil {
		httpClient.Transport = http.DefaultTransport
	}
	httpClient.Transport = &authTransport{
		underlying:  httpClient.Transport,
		accessToken: cfg.StravaAccessToken,
	}

	return &stravaApiHttpClient{
		httpClient:  httpClient,
		baseUrl:     cfg.StravaApiBaseUrl,
		accessToken: cfg.StravaAccessToken,
	}
}

// GetWattsStream implements [StravaApiClient].
func (c *stravaApiHttpClient) GetWattsStream(ctx context.Context, id int64) (*WattsStreamDto, error) {
	u := fmt.Sprintf("%s/activities/%d/streams?keys=watts&key_by_type=true", c.baseUrl, id)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

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

// GetDetailedActivityByID implements [StravaApiClient].
func (c *stravaApiHttpClient) GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error) {
	u := fmt.Sprintf("%s/activities/%d", c.baseUrl, activityID)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("strava error: %s", resp.Status)
	}

	var act DetailedActivityDto
	if err := json.NewDecoder(resp.Body).Decode(&act); err != nil {
		return nil, err
	}

	return &act, nil
}

// GetActivitiesByPeriod implements [StravaApiClient].
func (c *stravaApiHttpClient) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*ActivityDto, error) {
	startDate := period.StartDate()
	endDate := period.EndDate()

	start := getDate(startDate)
	end := getDate(endDate).Add(24 * time.Hour)
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+"/athlete/activities", nil)
	if err != nil {
		return nil, err
	}

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

func getDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
