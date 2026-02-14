package strava

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRefreshAccessToken_GivenStravaOauthHttpServerAndRefreshToken_WhenRefreshAccessToken_ThenItShouldCallCorrectEndpointAndDecodeResponse(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server", t, func() {
		request := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		var gotPath string
		var gotBody *RefreshAccessTokenRequest

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path

			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Errorf("failed to decode JSON: %v", err)
				return
			}

			_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"new-access-token","expires_at":1771095847,"expires_in":21600,"refresh_token":"test-refresh-token"}`))
		}))
		defer server.Close()
		client := NewStravaOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshAccessToken is called", func() {
			resp, err := client.RefreshAccessToken(request)

			Convey("Then it should call the correct endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/token")
			})

			Convey("Then it should include correct body", func() {
				So(gotBody, ShouldResemble, request)
			})

			Convey("Then it should decode the response correctly", func() {
				So(resp, ShouldNotBeNil)
				So(resp.TokenType, ShouldEqual, "Bearer")
				So(resp.AccessToken, ShouldEqual, "new-access-token")
				So(resp.ExpiresAt, ShouldEqual, 1771095847)
				So(resp.ExpiresIn, ShouldEqual, 21600)
				So(resp.RefreshToken, ShouldEqual, "test-refresh-token")
			})
		})
	})
}

func TestRefreshAccessToken_GivenStravaOauthHttpServerAndRefreshToken_WhenRefreshAccessTokenAndStravaReturnsError_ThenItShouldReturnError(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server that returns an error", t, func() {
		request := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}))
		defer server.Close()
		client := NewStravaOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshAccessToken is called", func() {
			resp, err := client.RefreshAccessToken(request)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestRefreshAccessToken_GivenStravaOauthHttpServerAndRefreshToken_WhenRefreshAccessTokenAndStravaReturnsInvalidJson_ThenItShouldReturnError(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server that returns invalid JSON", t, func() {
		request := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`invalid-json`))
		}))
		defer server.Close()
		client := NewStravaOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshAccessToken is called", func() {
			resp, err := client.RefreshAccessToken(request)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}
