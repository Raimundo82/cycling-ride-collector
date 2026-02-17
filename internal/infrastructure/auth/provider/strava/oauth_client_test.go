package strava

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	. "github.com/smartystreets/goconvey/convey"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRefreshTokenShouldCallCorrectEndpointAndDecodeResponseWithValidRequest(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server", t, func() {
		request := &provider.RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		var gotPath string
		var gotBody *provider.RefreshAccessTokenRequest

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path

			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Errorf("failed to decode JSON: %v", err)
				return
			}

			_, _ = w.Write([]byte(`{"token_type":"Bearer","access_token":"new-access-token","expires_at":1771095847,"expires_in":21600,"refresh_token":"test-refresh-token"}`))
		}))
		defer server.Close()
		client := NewOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

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

func TestRefreshTokenShouldReturnErrorWhenStravaReturnsError(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server that returns an error", t, func() {
		request := &provider.RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}))
		defer server.Close()
		client := NewOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestRefreshTokenShouldReturnErrorWhenStravaReturnsInvalidJson(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server that returns invalid JSON", t, func() {
		request := &provider.RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`invalid-json`))
		}))
		defer server.Close()
		client := NewOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "invalid character 'i' looking for beginning of value")
			})
		})
	})
}

func TestRefreshTokenShouldReturnErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and an invalid base URL", t, func() {
		request := &provider.RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		client := NewOAuthHttpClient(http.DefaultClient, &config.Config{StravaOauthBaseUrl: "://invalid-url"})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should return a request creation error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to create request")
			})
		})
	})
}

func TestRefreshTokenShouldReturnErrorWhenHTTPRequestFails(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and an http client transport error", t, func() {
		request := &provider.RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		httpErr := errors.New("transport error")
		httpClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return nil, httpErr
			}),
		}
		client := NewOAuthHttpClient(httpClient, &config.Config{StravaOauthBaseUrl: "http://example.com"})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should return the http client error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "transport error")
			})
		})
	})
}
