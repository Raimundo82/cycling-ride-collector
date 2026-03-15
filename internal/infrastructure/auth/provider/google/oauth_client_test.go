package auth_google

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestRefreshTokenShouldCallConfiguredEndpointAndSendFormEncodedBody(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a google oauth http server", t, func() {
		request := &auth_model.RefreshAccessTokenRequest{
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
			GrantType:    "refresh_token",
		}

		var gotPath string
		var gotContentType string
		var gotForm url.Values

		httpClient := &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				gotPath = r.URL.Path
				gotContentType = r.Header.Get("Content-Type")

				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				gotForm, err = url.ParseQuery(string(bodyBytes))
				if err != nil {
					t.Fatalf("failed to parse form: %v", err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader(`{"token_type":"Bearer","access_token":"new-google-access-token","expires_in":3600}`)),
					Header:     make(http.Header),
				}, nil
			}),
		}
		client := NewOAuthHttpClient(httpClient, &config.Config{
			GoogleOAuth: config.GoogleOAuthConfig{TokenURL: "http://example.com/token"},
		})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should call the configured endpoint", func() {
				So(err, ShouldBeNil)
				So(gotPath, ShouldEqual, "/token")
			})

			Convey("Then it should send an x-www-form-urlencoded body", func() {
				So(gotContentType, ShouldEqual, "application/x-www-form-urlencoded")
				So(gotForm.Get("client_id"), ShouldEqual, request.ClientID)
				So(gotForm.Get("client_secret"), ShouldEqual, request.ClientSecret)
				So(gotForm.Get("grant_type"), ShouldEqual, request.GrantType)
				So(gotForm.Get("refresh_token"), ShouldEqual, request.RefreshToken)
			})

			Convey("Then it should decode the response correctly", func() {
				So(resp, ShouldNotBeNil)
				So(resp.TokenType, ShouldEqual, "Bearer")
				So(resp.AccessToken, ShouldEqual, "new-google-access-token")
				So(resp.ExpiresIn, ShouldEqual, 3600)
			})
		})
	})
}

func TestRefreshTokenShouldUseDefaultGoogleTokenURLWhenConfigIsEmpty(t *testing.T) {
	Convey("Given an OAuth client created without a Google token URL", t, func() {
		client := NewOAuthHttpClient(http.DefaultClient, &config.Config{})

		Convey("Then it should use the default Google token URL", func() {
			So(client.baseUrl, ShouldEqual, defaultGoogleTokenURL)
		})
	})
}

func TestRefreshTokenShouldReturnErrorWhenGoogleReturnsError(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a google oauth http server that returns an error", t, func() {
		request := &auth_model.RefreshAccessTokenRequest{
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
			GrantType:    "refresh_token",
		}

		httpClient := &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadRequest,
					Status:     "400 Bad Request",
					Body:       io.NopCloser(strings.NewReader("Bad Request")),
					Header:     make(http.Header),
				}, nil
			}),
		}
		client := NewOAuthHttpClient(httpClient, &config.Config{
			GoogleOAuth: config.GoogleOAuthConfig{TokenURL: "http://example.com/token"},
		})

		Convey("When RefreshToken is called", func() {
			resp, err := client.RefreshToken(context.Background(), request)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestRefreshTokenShouldReturnErrorWhenGoogleReturnsInvalidJson(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a google oauth http server that returns invalid JSON", t, func() {
		request := &auth_model.RefreshAccessTokenRequest{
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
			GrantType:    "refresh_token",
		}

		httpClient := &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader(`invalid-json`)),
					Header:     make(http.Header),
				}, nil
			}),
		}
		client := NewOAuthHttpClient(httpClient, &config.Config{
			GoogleOAuth: config.GoogleOAuthConfig{TokenURL: "http://example.com/token"},
		})

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
		request := &auth_model.RefreshAccessTokenRequest{
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
			GrantType:    "refresh_token",
		}

		client := NewOAuthHttpClient(http.DefaultClient, &config.Config{
			GoogleOAuth: config.GoogleOAuthConfig{TokenURL: "://invalid-url"},
		})

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
		request := &auth_model.RefreshAccessTokenRequest{
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RefreshToken: "google-refresh-token",
			GrantType:    "refresh_token",
		}

		httpErr := errors.New("transport error")
		httpClient := &http.Client{
			Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
				return nil, httpErr
			}),
		}

		client := NewOAuthHttpClient(httpClient, &config.Config{
			GoogleOAuth: config.GoogleOAuthConfig{TokenURL: "http://example.com"},
		})

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
