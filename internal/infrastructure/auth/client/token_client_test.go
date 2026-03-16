package token_client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
	. "github.com/smartystreets/goconvey/convey"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestShouldCreateTokenClientWithTokenURLAndDefaultTimeoutWhenNewTokenClientIsCalled(t *testing.T) {
	Convey("Given a token url", t, func() {
		const tokenURL = "https://www.strava.com/oauth/token"

		Convey("When NewTokenClient is called", func() {
			client := NewTokenClient(tokenURL)

			Convey("Then it should create a token client with configured url and timeout", func() {
				So(client, ShouldNotBeNil)

				impl, ok := client.(*tokenClient)
				So(ok, ShouldBeTrue)
				So(impl, ShouldNotBeNil)
				So(impl.tokenUrl, ShouldEqual, tokenURL)
				So(impl.httpClient, ShouldNotBeNil)
				So(impl.httpClient.Timeout, ShouldEqual, 10*time.Second)
			})
		})
	})
}

func TestShouldReturnErrorWhenRequestCreationFailsDuringRefreshToken(t *testing.T) {
	Convey("Given a token client with invalid token url", t, func() {
		client := &tokenClient{tokenUrl: "://invalid-url", httpClient: &http.Client{Timeout: 10 * time.Second}}

		Convey("When RefreshToken is called", func() {
			output, err := client.RefreshToken(context.Background(), &token_model.RefreshTokenInput{})

			Convey("Then it should return request creation error", func() {
				So(output, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to create request")
			})
		})
	})
}

func TestShouldReturnErrorWhenDoRequestFailsDuringRefreshToken(t *testing.T) {
	Convey("Given a token client whose http client fails to do request", t, func() {
		expectedErr := errors.New("network unreachable")
		client := &tokenClient{
			tokenUrl: "https://www.strava.com/oauth/token",
			httpClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, expectedErr
			})},
		}

		Convey("When RefreshToken is called", func() {
			output, err := client.RefreshToken(context.Background(), &token_model.RefreshTokenInput{})

			Convey("Then it should return the do request error", func() {
				So(output, ShouldBeNil)
				So(errors.Is(err, expectedErr), ShouldBeTrue)
			})
		})
	})
}

func TestShouldReturnErrorWhenRefreshEndpointReturnsNonOKStatus(t *testing.T) {
	Convey("Given a token client and endpoint returning non ok status", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
		}))
		defer server.Close()

		client := &tokenClient{tokenUrl: server.URL, httpClient: server.Client()}

		Convey("When RefreshToken is called", func() {
			output, err := client.RefreshToken(context.Background(), &token_model.RefreshTokenInput{})

			Convey("Then it should return strava response status error", func() {
				So(output, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "error: 401 Unauthorized")
			})
		})
	})
}

func TestShouldReturnErrorWhenDecodeFailsForOKStatusResponseDuringRefreshToken(t *testing.T) {
	Convey("Given a token client and endpoint returning invalid json with status ok", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`invalid-json`))
		}))
		defer server.Close()

		client := &tokenClient{tokenUrl: server.URL, httpClient: server.Client()}

		Convey("When RefreshToken is called", func() {
			output, err := client.RefreshToken(context.Background(), &token_model.RefreshTokenInput{})

			Convey("Then it should return decode error", func() {
				So(output, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestShouldReturnAccessTokenWhenRefreshEndpointReturnsOKAndValidPayload(t *testing.T) {
	Convey("Given a token client and endpoint returning valid token payload", t, func() {
		var gotMethod string
		var gotContentType string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod = r.Method
			gotContentType = r.Header.Get("Content-Type")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"access_token":"new-access-token"}`))
		}))
		defer server.Close()

		client := &tokenClient{tokenUrl: server.URL, httpClient: server.Client()}
		input := &token_model.RefreshTokenInput{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			GrantType:    "refresh_token",
			RefreshToken: "refresh-token",
		}

		Convey("When RefreshToken is called", func() {
			output, err := client.RefreshToken(context.Background(), input)

			Convey("Then it should return decoded access token", func() {
				So(err, ShouldBeNil)
				So(output, ShouldNotBeNil)
				So(output.AccessToken, ShouldEqual, "new-access-token")
				So(gotMethod, ShouldEqual, http.MethodPost)
				So(gotContentType, ShouldEqual, "application/x-www-form-urlencoded")
			})
		})
	})
}
