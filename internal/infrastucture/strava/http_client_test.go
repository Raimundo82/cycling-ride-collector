package strava

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetActivitiesByDate_ConstructsCorrectRequestAndDecodesActivities(t *testing.T) {
	Convey("Given a strava http server and a date", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByDate is called", func() {
			acts, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/athlete/activities")
			})
			Convey("It should include day range", func() {
				start := date.Unix()
				end := date.Add(24 * time.Hour).Unix()

				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("after=%d", start))
				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("before=%d", end))
			})
			Convey("It should decode activities", func() {
				So(len(acts), ShouldEqual, 2)
			})
		})
	})
}

func TestGetActivitiesByDate_ReturnsEmptySliceWhenNoActivities(t *testing.T) {
	Convey("Given a strava http server and a date", t, func() {
		date := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivityByDate is called", func() {
			acts, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/athlete/activities")
			})
			Convey("It should include day range", func() {
				start := date.Unix()
				end := date.Add(24 * time.Hour).Unix()

				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("after=%d", start))
				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("before=%d", end))
			})
			Convey("It should decode activities", func() {
				So(len(acts), ShouldEqual, 0)
			})
		})
	})
}

func TestGetActivityByDate_ErrorCases(t *testing.T) {
	Convey("Given a strava http client", t, func() {
		client := NewHttpClient(http.DefaultClient, &config.Config{StravaApiBaseUrl: "http://invalid"})

		Convey("When context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			acts, err := client.GetActivitiesByDate(ctx, time.Now())

			Convey("It should return a context error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}

func TestGetActivitiesByDate_ReturnsErrorOnNonOKStatus(t *testing.T) {
	Convey("Given a strava http server that returns 401", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByDate is called", func() {
			acts, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should return an error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetActivitiesByDate_ReturnsErrorOnInvalidJSON(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByDate is called", func() {
			acts, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should return a decoding error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetActivitiesByDate_ReturnsErrorOnNetworkFailure(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		client := NewHttpClient(http.DefaultClient, &config.Config{StravaApiBaseUrl: "http://localhost:99999"})

		Convey("When GetActivitiesByDate is called", func() {
			acts, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should return a network error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetActivitiesByDate_ReturnsErrorOnTimeout(t *testing.T) {
	Convey("Given a slow strava http server", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`[{"id":1}]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByDate is called with a timeout context", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			acts, err := client.GetActivitiesByDate(ctx, date)

			Convey("It should return a timeout error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context deadline exceeded")
			})
		})
	})
}

func TestGetWattsStream_ConstructsCorrectRequestAndDecodesStream(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"watts":{"data": [103,100,50]}}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called", func() {
			activityID := int64(12345)
			watts, err := client.GetWattsStream(context.Background(), activityID)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, fmt.Sprintf("/activities/%d/streams", activityID))
			})

			Convey("It should include the correct query parameters", func() {
				So(gotQuery, ShouldEqual, "keys=watts&key_by_type=true")
			})

			Convey("It should decode the watts stream", func() {
				So(watts.WattsData, ShouldResemble, []int{103, 100, 50})
			})
		})
	})
}

func TestGetWattsStream_ReturnsEmptyWattsWhenNoWattsData(t *testing.T) {
	Convey("Given a strava http server that returns empty watts data", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"watts":{"data": null}}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called", func() {
			watts, err := client.GetWattsStream(context.Background(), activityID)

			Convey("It should return empty watts data", func() {
				So(watts, ShouldNotBeNil)
				So(watts.WattsData, ShouldBeEmpty)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetWattsStream_ReturnsEmptyWattsWhenWattsDataIsEmpty(t *testing.T) {
	Convey("Given a strava http server that returns empty watts data", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"watts":{"data": []}}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called", func() {
			watts, err := client.GetWattsStream(context.Background(), activityID)

			Convey("It should return empty watts data", func() {
				So(watts, ShouldNotBeNil)
				So(watts.WattsData, ShouldBeEmpty)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetWattsStream_ReturnsErrorOnNonOKStatus(t *testing.T) {
	Convey("Given a strava http server that returns 404", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called", func() {
			watts, err := client.GetWattsStream(context.Background(), activityID)

			Convey("It should return an error", func() {
				So(watts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetWattsStream_ReturnsErrorOnCanceledContext(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte(`{"watts":{"data": [100]}}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called with a canceled context", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			watts, err := client.GetWattsStream(ctx, 12345)

			Convey("It should return a context error", func() {
				So(watts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}

func TestGetWattsStream_ReturnsErrorOnInvalidJSON(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetWattsStream is called", func() {
			watts, err := client.GetWattsStream(context.Background(), 12345)

			Convey("It should return empty watts data", func() {
				So(watts, ShouldNotBeNil)
				So(watts.WattsData, ShouldBeEmpty)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetActivitiesByDate_SendsAuthorizationHeader(t *testing.T) {
	Convey("Given a strava http client with an access token", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		var gotAuthHeader string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuthHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`[{"id":1}]`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaApiBaseUrl:  server.URL,
			StravaAccessToken: "test_access_token_123",
		})

		Convey("When GetActivitiesByDate is called", func() {
			_, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should send the Authorization header", func() {
				So(err, ShouldBeNil)
				So(gotAuthHeader, ShouldEqual, "Bearer test_access_token_123")
			})
		})
	})
}

func TestGetActivitiesByDate_NoAuthHeaderWhenTokenEmpty(t *testing.T) {
	Convey("Given a strava http client without an access token", t, func() {
		date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		var gotAuthHeader string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuthHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`[{"id":1}]`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaApiBaseUrl:  server.URL,
			StravaAccessToken: "",
		})

		Convey("When GetActivitiesByDate is called", func() {
			_, err := client.GetActivitiesByDate(context.Background(), date)

			Convey("It should not send an Authorization header", func() {
				So(err, ShouldBeNil)
				So(gotAuthHeader, ShouldBeEmpty)
			})
		})
	})
}

func TestGetWattsStream_SendsAuthorizationHeader(t *testing.T) {
	Convey("Given a strava http client with an access token", t, func() {
		var gotAuthHeader string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuthHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{"watts":{"data": [100, 200]}}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaApiBaseUrl:  server.URL,
			StravaAccessToken: "test_access_token_456",
		})

		Convey("When GetWattsStream is called", func() {
			_, err := client.GetWattsStream(context.Background(), 12345)

			Convey("It should send the Authorization header", func() {
				So(err, ShouldBeNil)
				So(gotAuthHeader, ShouldEqual, "Bearer test_access_token_456")
			})
		})
	})
}

func TestRefreshAccessToken_SuccessfullyRefreshesToken(t *testing.T) {
	Convey("Given a strava http server that returns a successful token refresh", t, func() {
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{
				"access_token": "new_access_token_789",
				"refresh_token": "new_refresh_token_789",
				"expires_at": 1735257600,
				"expires_in": 21600
			}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaBaseUrl:      server.URL,
			StravaApiBaseUrl:   server.URL,
			StravaClientID:     "test_client_id",
			StravaClientSecret: "test_client_secret",
			StravaRefreshToken: "test_refresh_token",
		})

		Convey("When RefreshAccessToken is called", func() {
			tokenResp, err := client.RefreshAccessToken(context.Background())

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/oauth/token")
			})

			Convey("It should include the correct parameters", func() {
				So(gotQuery, ShouldContainSubstring, "client_id=test_client_id")
				So(gotQuery, ShouldContainSubstring, "client_secret=test_client_secret")
				So(gotQuery, ShouldContainSubstring, "grant_type=refresh_token")
				So(gotQuery, ShouldContainSubstring, "refresh_token=test_refresh_token")
			})

			Convey("It should decode the token response", func() {
				So(tokenResp, ShouldNotBeNil)
				So(tokenResp.AccessToken, ShouldEqual, "new_access_token_789")
				So(tokenResp.RefreshToken, ShouldEqual, "new_refresh_token_789")
				So(tokenResp.ExpiresAt, ShouldEqual, 1735257600)
				So(tokenResp.ExpiresIn, ShouldEqual, 21600)
			})

			Convey("It should update the client's access token", func() {
				So(client.accessToken, ShouldEqual, "new_access_token_789")
			})
		})
	})
}

func TestRefreshAccessToken_ReturnsErrorOnNonOKStatus(t *testing.T) {
	Convey("Given a strava http server that returns 401", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaBaseUrl:      server.URL,
			StravaApiBaseUrl:   server.URL,
			StravaClientID:     "test_client_id",
			StravaClientSecret: "test_client_secret",
			StravaRefreshToken: "test_refresh_token",
		})

		Convey("When RefreshAccessToken is called", func() {
			tokenResp, err := client.RefreshAccessToken(context.Background())

			Convey("It should return an error", func() {
				So(tokenResp, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava token refresh error")
			})
		})
	})
}

func TestRefreshAccessToken_ReturnsErrorOnInvalidJSON(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaBaseUrl:      server.URL,
			StravaApiBaseUrl:   server.URL,
			StravaClientID:     "test_client_id",
			StravaClientSecret: "test_client_secret",
			StravaRefreshToken: "test_refresh_token",
		})

		Convey("When RefreshAccessToken is called", func() {
			tokenResp, err := client.RefreshAccessToken(context.Background())

			Convey("It should return a decoding error", func() {
				So(tokenResp, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestRefreshAccessToken_ReturnsErrorOnCanceledContext(t *testing.T) {
	Convey("Given a strava http client", t, func() {
		client := NewHttpClient(http.DefaultClient, &config.Config{
			StravaBaseUrl:      "http://invalid",
			StravaApiBaseUrl:   "http://invalid",
			StravaClientID:     "test_client_id",
			StravaClientSecret: "test_client_secret",
			StravaRefreshToken: "test_refresh_token",
		})

		Convey("When context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			tokenResp, err := client.RefreshAccessToken(ctx)

			Convey("It should return a context error", func() {
				So(tokenResp, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}
