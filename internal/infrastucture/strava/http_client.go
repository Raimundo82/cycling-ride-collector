package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/config"
)

type httpClient struct {
	httpClient *http.Client
	baseUrl    string
}

var _ client = (*httpClient)(nil)

func NewHttpClient(http *http.Client, cfg *config.Config) *httpClient {
	return &httpClient{
		httpClient: http,
		baseUrl:    cfg.StravaApiBaseUrl,
	}
}

func (c *httpClient) GetActivitiesByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(time.Hour * 24)
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

func (c *httpClient) GetWattsStream(ctx context.Context, id int64) (*WattsStreamDto, error) {
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
