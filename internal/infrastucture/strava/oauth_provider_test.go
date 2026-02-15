package strava

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	. "github.com/smartystreets/goconvey/convey"
)

type stubOAuthClient struct {
	token *RefreshAccessTokenResponse
	err   error
}

// RefreshAccessToken implements [StravaOAuthClient].
func (s *stubOAuthClient) RefreshAccessToken(req *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error) {
	return s.token, s.err
}

var _ StravaOAuthClient = (*stubOAuthClient)(nil)

func TestStravaOAuthHttpClient_GivenStravaOAuthClient_WhenRequestIsValid_ThenItShouldReturnNewAccessToken(t *testing.T) {
	Convey("Given a Strava OAuth HTTP client", t, func() {
		refreshTokenResponse := &RefreshAccessTokenResponse{
			TokenType:   "Bearer",
			AccessToken: "new-access-token",
			ExpiresAt:   1771095847,
			ExpiresIn:   21600,
		}

		refreshTokenRequest := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}
		provider := &StravaOAuthProvider{oauthClient: &stubOAuthClient{token: refreshTokenResponse}}
		client := provider.oauthClient

		Convey("When the request is valid", func() {
			resp, err := client.RefreshAccessToken(refreshTokenRequest)

			Convey("Then it should return the new access token", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp, ShouldResemble, refreshTokenResponse)
			})
		})
	})
}

func TestStravaOAuthHttpClient_GivenStravaOAuthClient_WhenRequestIsInvalid_ThenItShouldReturnError(t *testing.T) {
	Convey("Given a Strava OAuth HTTP client that returns an error", t, func() {
		provider := &StravaOAuthProvider{oauthClient: &stubOAuthClient{err: http.ErrServerClosed}}
		client := provider.oauthClient

		refreshTokenRequest := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		Convey("When the request is invalid", func() {
			resp, err := client.RefreshAccessToken(refreshTokenRequest)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(resp, ShouldBeNil)
			})
		})
	})
}

func TestRefreshAccessToken_GivenStravaOauthHttpServerAndRefreshToken_WhenRefreshAccessTokenAndStravaReturnsValidResponse_ThenItShouldReturnNewAccessToken(t *testing.T) {
	Convey("Given a RefreshAccessTokenRequest and a strava oauth http server that returns a valid response", t, func() {
		request := &RefreshAccessTokenRequest{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
			GrantType:    "refresh_token",
		}

		response := &RefreshAccessTokenResponse{
			TokenType:    "Bearer",
			AccessToken:  "new-access-token",
			ExpiresAt:    1771095847,
			ExpiresIn:    21600,
			RefreshToken: "test-refresh-token",
		}
		gotBody := &RefreshAccessTokenRequest{}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
			_ = json.NewDecoder(r.Body).Decode(gotBody)
		}))

		defer server.Close()

		provider := &StravaOAuthProvider{oauthClient: NewStravaOAuthHttpClient(server.Client(), &config.Config{StravaOauthBaseUrl: server.URL})}
		client := provider.oauthClient

		Convey("When RefreshAccessToken is called", func() {
			resp, _ := client.RefreshAccessToken(request)

			Convey("Then it should send the correct request to Strava", func() {
				So(gotBody, ShouldNotBeNil)
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
