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

func TestGetActivityByDate_BuildsCorrectRequestAndDecodes(t *testing.T) {
	Convey("Given a strava http server", t, func() {
		var gotPath string
		var gotQuery string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			_, _ = w.Write([]byte(`[{"id":1},{"id":2}]`))
		}))
		defer server.Close()
		client := NewHttpClient(server.Client(), &config.Config{StravaApiBaseUrl: server.URL})
		Convey("When GetActivityByDate is called", func() {
			date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
			acts, err := client.GetActivityByDate(context.Background(), date)

			Convey("It should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/athlete/activities")
			})
			Convey("It should include day range", func() {
				start := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC).Unix()
				end := start + 86400

				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("after=%d", start))
				So(gotQuery, ShouldContainSubstring, fmt.Sprintf("before=%d", end))
			})
			Convey("It should decode activities", func() {
				So(len(acts), ShouldEqual, 2)
			})
		})
	})
}

func TestGetWattsStream(t *testing.T) {
	Convey("Given", t, func() {
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
				So(gotQuery, ShouldEqual, "keys=watts&key_by_type=true")
			})

			Convey("Then", func() {
				So(watts, ShouldResemble, &WattsStreamDto{WattsData: []int{103, 100, 50}})
			})
		})
	})
}
