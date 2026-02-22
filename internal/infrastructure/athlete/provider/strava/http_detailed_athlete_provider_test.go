package athlete_strava

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreatingHttpAthleteStatsProviderReturnsANonNilHttpAthleteStatsProvider(t *testing.T) {
	Convey("Given a valid http client and base url", t, func() {
		httpClient := &http.Client{}
		baseUrl := "https://www.strava.com/api/v3"

		Convey("When creating a new http athlete stats provider", func() {
			provider := NewHttpAthleteStatsProvider(httpClient, baseUrl)

			Convey("Then the provider should be created successfully", func() {
				So(provider, ShouldNotBeNil)
				So(provider.baseUrl, ShouldEqual, baseUrl)
				So(provider.httpClient, ShouldResemble, *httpClient)
			})
		})
	})
}

func TestGetAthleteDataReturnsDetailedAthleteData(t *testing.T) {
	Convey("Given a http athlete stats provider and a mock strava server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": 12345,
				"weight": 70.5
			}`))
		}))
		defer server.Close()

		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When calling GetAthleteData", func() {
			athleteData, err := provider.GetDetailedAthlete(context.Background())

			Convey("Then it should return a detailed athlete data", func() {
				So(athleteData, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(athleteData.ID, ShouldEqual, 12345)
				So(athleteData.Weight, ShouldEqual, 70.5)
			})
		})
	})
}

func TestGetAthleteDataReturnsErrorWhenContextIsCanceled(t *testing.T) {
	Convey("Given a http athlete stats provider with an invalid URL", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "http://invalid-url")

		Convey("When context is canceled", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			athleteData, err := provider.GetDetailedAthlete(ctx)

			Convey("It should return a context error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context canceled")
			})
		})
	})
}

func TestGetAthleteDataReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a http athlete stats provider with an invalid URL", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "://invalid-url")

		Convey("When GetDetailedAthlete is called", func() {
			athleteData, err := provider.GetDetailedAthlete(context.Background())

			Convey("Then it should return request creation error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "missing protocol scheme")
			})
		})
	})
}

func TestGetAthleteDataReturnsErrorWhenStatusIsNonOK(t *testing.T) {
	Convey("Given a http athlete stats provider and a mock strava server that returns 401", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetDetailedAthlete is called", func() {
			athleteData, err := provider.GetDetailedAthlete(context.Background())

			Convey("It should return an error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "strava error")
			})
		})
	})
}

func TestGetAthleteDataReturnsDecodingErrorWhenJSONIsInvalid(t *testing.T) {
	Convey("Given a strava http server that returns invalid JSON", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{invalid json}`))
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetDetailedAthlete is called", func() {
			athleteData, err := provider.GetDetailedAthlete(context.Background())

			Convey("It should return a decoding error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid character")
			})
		})
	})
}

func TestGetDetailedAthleteReturnsErrorWhenNetworkFails(t *testing.T) {
	Convey("Given an unreachable strava server", t, func() {
		provider := NewHttpAthleteStatsProvider(http.DefaultClient, "http://localhost:99999")

		Convey("When GetDetailedAthlete is called", func() {
			athleteData, err := provider.GetDetailedAthlete(context.Background())

			Convey("It should return a network error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid port")
			})
		})
	})
}

func TestGetDetailedAthleteReturnsErrorWhenRequestTimeouts(t *testing.T) {
	Convey("Given a slow strava http server", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			_, _ = w.Write([]byte(`{"id":1}`))
		}))
		defer server.Close()
		provider := NewHttpAthleteStatsProvider(server.Client(), server.URL)

		Convey("When GetDetailedAthlete is called with a timeout context", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			athleteData, err := provider.GetDetailedAthlete(ctx)

			Convey("It should return a timeout error", func() {
				So(athleteData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "context deadline exceeded")
			})
		})
	})
}
