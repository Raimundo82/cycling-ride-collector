package activity_strava

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/domain"
	. "github.com/smartystreets/goconvey/convey"
)

const TOKEN = "test-access-token"

func TestGetActivitiesByPeriodCallsCorrectEndpointAndDecodesActivities(t *testing.T) {
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

		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetActivitiesByPeriodReturnsEmptySliceWhenNoActivitiesInPeriod(t *testing.T) {
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
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetActivitiesByPeriodReturnsContextErrorWhenContextIsCanceled(t *testing.T) {
	Convey("Given a strava http client", t, func() {
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "http://invalid"}})
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

func TestGetActivitiesByPeriodReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a strava activity provider with invalid base URL", t, func() {
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "://invalid-url"}})
		startDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("Then it should return request creation error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing protocol scheme")
			})
		})
	})
}

func TestGetActivitiesByPeriodReturnsErrorWhenStatusIsNonOK(t *testing.T) {
	Convey("Given a strava http server that returns 401", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetActivitiesByPeriodReturnsDecodingErrorWhenJSONIsInvalid(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should return a decoding error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid character")
			})
		})
	})
}

func TestGetActivitiesByPeriodReturnsErrorWhenNetworkFails(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "http://localhost:99999"}})
		Convey("When GetActivitiesByPeriod is called", func() {
			acts, err := client.GetActivitiesByPeriod(context.Background(), period)

			Convey("It should return a network error", func() {
				So(acts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid port")
			})
		})
	})
}

func TestGetActivitiesByPeriodReturnsErrorWhenRequestTimeouts(t *testing.T) {
	Convey("Given a slow strava http server", t, func() {
		startDate := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		endtDate := time.Date(2024, 1, 11, 0, 0, 0, 0, time.UTC)
		period, _ := domain.NewPeriod(startDate, endtDate)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`[{"id":1}]`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamConstructsCorrectRequestAndDecodesStream(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`{"watts":{"data": [103,100,50]}}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsEmptyWattsWhenNoWattsData(t *testing.T) {
	Convey("Given a strava http server that returns empty watts data", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"watts":{"data": null}}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsEmptyWattsWhenWattsDataIsEmpty(t *testing.T) {
	Convey("Given a strava http server that returns empty watts data", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"watts":{"data": []}}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsErrorOnNonOKStatus(t *testing.T) {
	Convey("Given a strava http server that returns 404", t, func() {
		activityID := int64(12345)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsErrorOnCanceledContext(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte(`{"watts":{"data": [100]}}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsErrorOnInvalidJSON(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

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

func TestGetWattsStreamReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a strava activity provider with invalid base URL", t, func() {
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "://invalid-url"}})

		Convey("When GetWattsStream is called", func() {
			watts, err := client.GetWattsStream(context.Background(), 12345)

			Convey("Then it should return request creation error", func() {
				So(watts, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing protocol scheme")
			})
		})
	})
}

func TestGetDetailedActivityByIDConstructsCorrectRequestAndDecodesActivity(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		var gotPath string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(`{"private_note": "Boas"}`))
		}))
		defer server.Close()

		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

		Convey("When GetDetailedActivityByID is called", func() {
			activityID := int64(98765)
			data, err := client.GetDetailedActivityByID(context.Background(), activityID)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, fmt.Sprintf("/activities/%d", activityID))
			})

			Convey("It should decode detailed activity data", func() {
				So(data.LegSensations, ShouldEqual, "Boas")
			})
		})
	})
}

func TestGetDetailedActivityByIDReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a strava activity provider with invalid base URL", t, func() {
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "://invalid-url"}})

		Convey("When GetDetailedActivityByID is called", func() {
			data, err := client.GetDetailedActivityByID(context.Background(), 123)

			Convey("Then it should return request creation error", func() {
				So(data, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing protocol scheme")
			})
		})
	})
}

func TestGetDetailedActivityByIDReturnsErrorWhenNetworkFails(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		client := NewActivityProvider(http.DefaultClient, &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: "http://localhost:99999"}})

		Convey("When GetDetailedActivityByID is called", func() {
			data, err := client.GetDetailedActivityByID(context.Background(), 123)

			Convey("Then it should return transport error", func() {
				So(data, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid port")
			})
		})
	})
}

func TestGetDetailedActivityByIDReturnsErrorWhenStatusIsNonOK(t *testing.T) {
	Convey("Given a strava http server that returns 401", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

		Convey("When GetDetailedActivityByID is called", func() {
			data, err := client.GetDetailedActivityByID(context.Background(), 123)

			Convey("Then it should return strava error", func() {
				So(data, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetDetailedActivityByIDReturnsErrorWhenJSONIsInvalid(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()

		client := NewActivityProvider(server.Client(), &config.Config{Strava: &config.StravaConfig{ApiBaseUrl: server.URL}})

		Convey("When GetDetailedActivityByID is called", func() {
			data, err := client.GetDetailedActivityByID(context.Background(), 123)

			Convey("Then it should return decoding error", func() {
				So(data, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid character")
			})
		})
	})
}
