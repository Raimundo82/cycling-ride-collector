package athlete_strava

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAthleteZonesReturnsAthleteZones(t *testing.T) {
	Convey("Given a http athlete stats provider and a mock strava server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{ 
				"power": { "zones": [{ "max": 200, "min": 100 }, { "max": 300, "min": 200 }] },
				"heart_rate": { "zones": [{ "max": 120, "min": 100 }, { "max": 150, "min": 120 }] } 
			}`))
		}))
		defer server.Close()

		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When calling GetAthleteZones", func() {
			athleteZones, err := provider.GetAthleteZones(context.Background())

			Convey("Then it should return athlete zones", func() {
				So(athleteZones, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(athleteZones.PowerRangeZones.Zones[0].Max, ShouldEqual, 200)
				So(athleteZones.PowerRangeZones.Zones[0].Min, ShouldEqual, 100)
				So(athleteZones.HeartRateRangeZones.Zones[0].Max, ShouldEqual, 120)
				So(athleteZones.HeartRateRangeZones.Zones[0].Min, ShouldEqual, 100)
			})
		})
	})
}

func TestGetAthleteZonesReturnsErrorWhenContextIsCanceled(t *testing.T) {
	Convey("Given a http athlete stats provider with an invalid URL", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "http://invalid-url")

		Convey("When context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			athleteZones, err := provider.GetAthleteZones(ctx)

			Convey("It should return a context error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}

func TestGetAthleteZonesReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a http athlete stats provider with an invalid URL", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "://invalid-url")

		Convey("When GetAthleteZones is called", func() {
			athleteZones, err := provider.GetAthleteZones(context.Background())

			Convey("Then it should return request creation error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing protocol scheme")
			})
		})
	})
}

func TestGetAthleteZonesReturnsErrorWhenStatusIsNonOK(t *testing.T) {
	Convey("Given a http athlete stats provider and a mock strava server that returns 401", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetAthleteZones is called", func() {
			athleteZones, err := provider.GetAthleteZones(context.Background())

			Convey("It should return an error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetAthleteZonesReturnsDecodingErrorWhenJSONIsInvalid(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetAthleteZones is called", func() {
			athleteZones, err := provider.GetAthleteZones(context.Background())

			Convey("It should return a decoding error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid character")
			})
		})
	})
}

func TestGetAthleteZonesReturnsErrorWhenNetworkFails(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "http://localhost:99999")

		Convey("When GetAthleteZones is called", func() {
			athleteZones, err := provider.GetAthleteZones(context.Background())

			Convey("It should return a network error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid port")
			})
		})
	})
}

func TestGetAthleteZonesReturnsErrorWhenRequestTimeouts(t *testing.T) {
	Convey("Given a slow strava http server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`{"id":1}`))
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetAthleteZones is called with a timeout context", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			athleteZones, err := provider.GetAthleteZones(ctx)

			Convey("It should return a timeout error", func() {
				So(athleteZones, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context deadline exceeded")
			})
		})
	})
}
