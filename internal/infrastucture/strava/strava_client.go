package strava

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type client interface {
	GetActivityByDate(ctx context.Context, date time.Time) ([]*ActivityDto, error)
	GetWattsStream(ctx context.Context, activityID int64) (*WattsStreamDto, error)
}

type StravaClient struct {
	httpClient   *http.Client
	oauth2Config *oauth2.Config
	tokenSource  oauth2.TokenSource
}

type StravaConfig struct {
	CLientUrl    string
	ClientApiUrl string
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
}

func NewStravaClient(config StravaConfig) *StravaClient {
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.CLientUrl + "/oauth/authorize",
			TokenURL: config.CLientUrl + "/oauth/token",
		},
	}

	token := &oauth2.Token{
		AccessToken:  config.AccessToken,
		RefreshToken: config.RefreshToken,
		TokenType:    "Bearer",
	}

	tokenSource := oauth2Config.TokenSource(context.Background(), token)
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	return &StravaClient{
		httpClient:   httpClient,
		oauth2Config: oauth2Config,
		tokenSource:  tokenSource,
	}
}
