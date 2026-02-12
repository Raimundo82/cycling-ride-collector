package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/config"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

type stravaHttpClient struct {
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
	_ StravaClient      = (*stravaHttpClient)(nil)
	_ http.RoundTripper = (*authTransport)(nil)
)

func NewHttpClient(httpClient *http.Client, cfg *config.Config) *stravaHttpClient {
	if httpClient.Transport == nil {
		httpClient.Transport = http.DefaultTransport
	}
	httpClient.Transport = &authTransport{
		underlying:  httpClient.Transport,
		accessToken: cfg.StravaAccessToken,
	}

	return &stravaHttpClient{
		httpClient:  httpClient,
		baseUrl:     cfg.StravaApiBaseUrl,
		accessToken: cfg.StravaAccessToken,
	}
}

// GetWattsStream implements [StravaClient].
func (c *stravaHttpClient) GetWattsStream(ctx context.Context, id int64) (*WattsStreamDto, error) {
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

// GetDetailedActivityByID implements [StravaClient].
func (c *stravaHttpClient) GetDetailedActivityByID(ctx context.Context, activityID int64) (*DetailedActivityDto, error) {
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

// GetActivitiesByPeriod implements [StravaClient].
func (c *stravaHttpClient) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*ActivityDto, error) {
	startDate := period.StartDate()
	endDate := period.EndDate()

	start := getDate(startDate)
	end := getDate(endDate)
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
