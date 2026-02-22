package activity_strava

import (
	"errors"
	"net/http"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/interfaces"
	. "github.com/smartystreets/goconvey/convey"
)

type authTransportMockTokenProvider struct {
	token string
	err   error
}

func (m *authTransportMockTokenProvider) GetValidToken() (string, error) {
	return m.token, m.err
}

type authTransportRoundTripFunc func(*http.Request) (*http.Response, error)

func (f authTransportRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var _ activity_interfaces.TokenProvider = (*authTransportMockTokenProvider)(nil)

const uri = "http://example.com"

func TestRoundTripReturnsWrappedErrorWhenTokenProviderFails(t *testing.T) {
	Convey("Given an auth transport whose token provider fails", t, func() {
		tokenErr := errors.New("token provider failed")
		underlyingCalled := false
		transport := &authTransport{
			underlying: authTransportRoundTripFunc(func(*http.Request) (*http.Response, error) {
				underlyingCalled = true
				return &http.Response{StatusCode: http.StatusOK}, nil
			}),
			tokenProvider: &authTransportMockTokenProvider{err: tokenErr},
		}
		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		Convey("When RoundTrip is called", func() {
			resp, err := transport.RoundTrip(req)

			Convey("Then it should return a wrapped token retrieval error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to get valid token for request")
				So(errors.Is(err, tokenErr), ShouldBeTrue)
				So(underlyingCalled, ShouldBeFalse)
			})
		})
	})
}

func TestRoundTripAddsAuthorizationHeaderWhenTokenIsAvailable(t *testing.T) {
	Convey("Given an auth transport with a valid access token", t, func() {
		const token = "abc123"

		var gotReq *http.Request
		transport := &authTransport{
			underlying: authTransportRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				gotReq = req
				return &http.Response{StatusCode: http.StatusAccepted}, nil
			}),
			tokenProvider: &authTransportMockTokenProvider{token: token},
		}

		originalReq, _ := http.NewRequest(http.MethodGet, uri, nil)
		originalReq.Header.Set("X-Test", "value")

		Convey("When RoundTrip is called", func() {
			resp, err := transport.RoundTrip(originalReq)

			Convey("Then it should clone request and add bearer token only to forwarded request", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusAccepted)
				So(gotReq, ShouldNotBeNil)
				So(gotReq, ShouldNotEqual, originalReq)
				So(gotReq.Header.Get("Authorization"), ShouldEqual, "Bearer "+token)
				So(originalReq.Header.Get("Authorization"), ShouldBeEmpty)
				So(gotReq.Header.Get("X-Test"), ShouldEqual, "value")
			})
		})
	})
}

func TestRoundTripDoesNotAddAuthorizationHeaderWhenTokenIsEmpty(t *testing.T) {
	Convey("Given an auth transport with empty token", t, func() {
		var gotReq *http.Request
		transport := &authTransport{
			underlying: authTransportRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				gotReq = req
				return &http.Response{StatusCode: http.StatusNoContent}, nil
			}),
			tokenProvider: &authTransportMockTokenProvider{token: ""},
		}

		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		Convey("When RoundTrip is called", func() {
			resp, err := transport.RoundTrip(req)

			Convey("Then it should pass through original request without authorization header", func() {
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
				So(gotReq, ShouldEqual, req)
				So(gotReq.Header.Get("Authorization"), ShouldBeEmpty)
			})
		})
	})
}

func TestRoundTripReturnsUnderlyingTransportError(t *testing.T) {
	Convey("Given an auth transport with valid token and failing underlying transport", t, func() {
		transportErr := errors.New("round trip failed")
		transport := &authTransport{
			underlying: authTransportRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				So(req.Header.Get("Authorization"), ShouldEqual, "Bearer token")
				return nil, transportErr
			}),
			tokenProvider: &authTransportMockTokenProvider{token: "token"},
		}

		req, _ := http.NewRequest(http.MethodGet, uri, nil)

		Convey("When RoundTrip is called", func() {
			resp, err := transport.RoundTrip(req)

			Convey("Then it should return the underlying transport error", func() {
				So(resp, ShouldBeNil)
				So(err, ShouldEqual, transportErr)
			})
		})
	})
}
