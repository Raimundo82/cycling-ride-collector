package activity_strava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	activity_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/model"
)

type activityProvider struct {
	httpClient *http.Client
	baseUrl    string
}

var _ ActivityProvider = (*activityProvider)(nil)

const stravaError = "strava error: %s"

func NewActivityProvider(httpClient *http.Client, cfg *config.Config) *activityProvider {
	return &activityProvider{
		httpClient: httpClient,
		baseUrl:    cfg.StravaApiBaseUrl,
	}
}

// GetActivitiesByPeriod implements [ActivityProvider].
func (a *activityProvider) GetActivitiesByPeriod(ctx context.Context, period domain.Period) ([]*activity_model.ActivityDto, error) {
	start := getDate(period.StartDate())
	end := getDate(period.EndDate()).Add(24 * time.Hour)

	req, err := http.NewRequestWithContext(ctx, "GET", a.baseUrl+"/athlete/activities", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("after", fmt.Sprint(start.Unix()))
	q.Set("before", fmt.Sprint(end.Unix()))
	req.URL.RawQuery = q.Encode()

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(stravaError, resp.Status)
	}

	var periodActivities []*activity_model.ActivityDto
	if err := json.NewDecoder(resp.Body).Decode(&periodActivities); err != nil {
		return nil, err
	}

	return periodActivities, nil
}

// GetDetailedActivityByID implements [ActivityProvider].
func (a *activityProvider) GetDetailedActivityByID(ctx context.Context, activityID int64) (*activity_model.DetailedActivityDto, error) {
	url := fmt.Sprintf("%s/activities/%d", a.baseUrl, activityID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(stravaError, resp.Status)
	}

	var activity activity_model.DetailedActivityDto
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return nil, err
	}

	return &activity, nil
}

// GetWattsStream implements [ActivityProvider].
func (a *activityProvider) GetWattsStream(ctx context.Context, activityID int64) (*activity_model.WattsStreamDto, error) {
	u := fmt.Sprintf("%s/activities/%d/streams?keys=watts&key_by_type=true", a.baseUrl, activityID)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(stravaError, resp.Status)
	}

	var streams activity_model.WattsStreamResponse

	err = json.NewDecoder(resp.Body).Decode(&streams)
	if err != nil {
		return &activity_model.WattsStreamDto{WattsData: []int{}}, nil
	}

	if len(streams.Watts.WattsData) == 0 {
		return &activity_model.WattsStreamDto{WattsData: []int{}}, nil
	}

	return &streams.Watts, nil
}

func getDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
