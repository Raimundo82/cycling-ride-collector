package strava

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/config"
	"github.com/raimundo82/go-strava-weekly/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetActivitiesByPeriod_GivenStravaHttpServerAndPeriod_WhenGetActivitiesByPeriod_ThenItShouldCallCorrectEndpointAndDecodeActivities(t *testing.T) {
	Convey("Given a strava http server that returns two activities for a period", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("Then it should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/athlete/activities")
			})
			Convey("Then it should include day range", func() {
				start := startDate.Unix()
				end := endtDate.Add(24 * time.Hour).Unix()

				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("after=%d", start))
				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("before=%d", end))
			})
			Convey("Then it should decode activities", func() {
				So(len(acts), ShouldEqual, 2)
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenNoActivitiesForPeriod_WhenGetActivitiesByPeriod_ThenItShouldReturnEmptySlice(t *testing.T) {
	Convey("Given a strava http server that returns no activities for a period", t, func() {
		startDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/athlete/activities")
			})
			Convey("It should include day range", func() {
				start := startDate.Unix()
				end := endtDate.Add(24 * time.Hour).Unix()

				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("after=%d", start))
				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("before=%d", end))
			})
			Convey("It should decode activities", func() {
				So(len(acts), ShouldEqual, 0)
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenContextIsCanceled_ThenItShouldReturnContextError(t *testing.T) {
	Convey("Given a strava http client", t, func() {
		client := NewHttpClient(http.DefaultClient, &config.Config{StravaApiBaseUrl: "http://invalid"})
		startDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		Convey("When context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			acts, err := client.GetActivitiesByPeriod(ctx, period)

			Convey("It should return a context error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenNonOKStatus_ThenItShouldReturnError(t *testing.T) {
	Convey("Given a strava http server that returns 401", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should return an error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenInvalidJSON_ThenItShouldReturnDecodingError(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should return a decoding error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenNetworkFailure_ThenItShouldReturnError(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)
		client := NewHttpClient(http.DefaultClient, &config.Config{StravaApiBaseUrl: "http://localhost:99999"})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should return a network error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenTimeout_ThenItShouldReturnError(t *testing.T) {
	Convey("Given a slow strava http server", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`[{"id":1}]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})

		Convey("When GetActivitiesByPeriod is called with a timeout context", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			acts, err := client.GetActivitiesByPeriod(ctx, period)

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

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenSendsAuthorizationHeader_ThenItShouldSendAuthorizationHeader(t *testing.T) {
	Convey("Given a strava http client with an access token", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endDate)
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

		Convey("When GetActivitiesByPeriod is called", func() {
			_, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should send the Authorization header", func() {
				So(err, ShouldBeNil)
				So(gotAuthHeader, ShouldEqual, "Bearer test_access_token_123")
			})
		})
	})
}

func TestGetActivitiesByPeriod_GivenStravaHttpClient_WhenNoAuthHeaderWhenTokenEmpty_ThenItShouldNotSendAuthorizationHeader(t *testing.T) {
	Convey("Given a strava http client without an access token", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endDate)
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

		Convey("When GetActivitiesByPeriod is called", func() {
			_, err := client.GetActivitiesByPeriod(context.Background(), period)

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

func TestGetWattsStream_NoAuthHeaderWhenTokenEmpty(t *testing.T) {
	Convey("Given a strava http client without an access token", t, func() {
		var gotAuthHeader string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuthHeader = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{"watts":{"data": [100, 200]}}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaApiBaseUrl:  server.URL,
			StravaAccessToken: "",
		})

		Convey("When GetWattsStream is called", func() {
			_, err := client.GetWattsStream(context.Background(), 12345)

			Convey("It should not send an Authorization header", func() {
				So(err, ShouldBeNil)
				So(gotAuthHeader, ShouldBeEmpty)
			})
		})
	})
}

func TestGetDetailedActivityByID_GivenStravaHttpClient_WhenGetsDetailedActivityByID_ThenItShouldConstructCorrectRequestAndDecodeActivity(t *testing.T) {
	Convey("Given a strava http server and a token", t, func() {
		var gotAuthHeader string
		var gotPath string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotAuthHeader = r.Header.Get("Authorization")
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"private_note": "Boas"}`))
		}))
		defer server.Close()

		client := NewHttpClient(server.Client(), &config.Config{
			StravaApiBaseUrl:  server.URL,
			StravaAccessToken: "test_access_token_789",
		})

		Convey("When GetDetailedActivityByID is called", func() {
			activityID := int64(98765)
			data, err := client.GetDetailedActivityByID(context.Background(), activityID)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, fmt.Sprintf("/activities/%d", activityID))
			})

			Convey("It should decode detailed activity data", func() {
				So(gotAuthHeader, ShouldEqual, "Bearer test_access_token_789")
				So(data.LegSensations, ShouldEqual, "Boas")
			})
		})
	})
}
