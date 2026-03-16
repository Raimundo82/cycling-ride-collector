package email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/raimundo82/cycling-ride-collector/internal/config"
	token_provider "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/provider"
	. "github.com/smartystreets/goconvey/convey"
)

type stubGoogleTokenProvider struct {
	Token string
	Err   error
	Ctx   context.Context
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

var _ token_provider.TokenProvider = (*stubGoogleTokenProvider)(nil)

func (s *stubGoogleTokenProvider) GetValidToken(ctx context.Context) (string, error) {
	s.Ctx = ctx
	return s.Token, s.Err
}

func TestNewGmailWorkoutReportSenderShouldConfigureSenderFromConfig(t *testing.T) {
	Convey("Given an http client, token provider, and config", t, func() {
		httpClient := &http.Client{}
		tokenProvider := &stubGoogleTokenProvider{}
		cfg := &config.Config{Email: &config.EmailConfig{From: "from@example.com", To: "to@example.com", Subject: "Subject"}}

		Convey("When NewGmailWorkoutReportSender is called", func() {
			sender := NewGmailWorkoutReportSender(httpClient, tokenProvider, cfg)

			Convey("Then it should map dependencies and email settings", func() {
				So(sender, ShouldNotBeNil)
				So(sender.httpClient, ShouldEqual, httpClient)
				So(sender.tokenProvider, ShouldEqual, tokenProvider)
				So(sender.from, ShouldEqual, "from@example.com")
				So(sender.to, ShouldEqual, "to@example.com")
				So(sender.subject, ShouldEqual, "Subject")
				So(sender.apiURL, ShouldEqual, gmailSendMessageURL)
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendSuccess(t *testing.T) {
	Convey("Given a gmail sender and a valid report", t, func() {
		var authorizationHeader string
		var rawMessage string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader = r.Header.Get("Authorization")
			var payload gmailSendRequest
			_ = json.NewDecoder(r.Body).Decode(&payload)
			rawMessage = payload.Raw
			_, _ = w.Write([]byte(`{"id":"gmail-message-id"}`))
		}))
		defer server.Close()

		reportPath := createTempReportFile(t)
		tokenProvider := &stubGoogleTokenProvider{Token: "google-access-token"}
		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: tokenProvider,
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        server.URL,
		}

		Convey("When Send is called", func() {
			err := sender.Send(reportPath)

			Convey("Then it sends the Gmail API request with the encoded message", func() {
				So(err, ShouldBeNil)
				So(authorizationHeader, ShouldEqual, "Bearer google-access-token")
				So(tokenProvider.Ctx, ShouldNotBeNil)

				decodedMessage, decodeErr := base64.RawURLEncoding.DecodeString(rawMessage)
				So(decodeErr, ShouldBeNil)
				So(string(decodedMessage), ShouldContainSubstring, "Subject: Cycling Workout Report")
				So(string(decodedMessage), ShouldContainSubstring, "filename=\"report.xlsx\"")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendCreatesRequestWithDeadlineContext(t *testing.T) {
	Convey("Given a gmail sender with a transport spy", t, func() {
		var requestHasDeadline bool
		httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			_, requestHasDeadline = req.Context().Deadline()
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"id":"gmail-message-id"}`)),
				Header:     make(http.Header),
			}, nil
		})}

		sender := &gmailWorkoutReportSender{
			httpClient:    httpClient,
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        "https://gmail.googleapis.com/gmail/v1/users/me/messages/send",
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should create request with deadline context", func() {
				So(err, ShouldBeNil)
				So(requestHasDeadline, ShouldBeTrue)
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendPassesCancellableContextToTokenProvider(t *testing.T) {
	Convey("Given a gmail sender with token provider spy", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"id":"gmail-message-id"}`))
		}))
		defer server.Close()

		tokenProvider := &stubGoogleTokenProvider{Token: "google-access-token"}
		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: tokenProvider,
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        server.URL,
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should pass a context with deadline to token provider", func() {
				So(err, ShouldBeNil)
				So(tokenProvider.Ctx, ShouldNotBeNil)
				deadline, ok := tokenProvider.Ctx.Deadline()
				So(ok, ShouldBeTrue)
				So(deadline.IsZero(), ShouldBeFalse)
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorWhenTokenProviderFails(t *testing.T) {
	Convey("Given a gmail sender with failing token provider", t, func() {
		sender := &gmailWorkoutReportSender{
			httpClient:    &http.Client{},
			tokenProvider: &stubGoogleTokenProvider{Err: errors.New("token unavailable")},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        "http://localhost",
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should return token retrieval error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to get google access token")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorWhenReportFileCannotBeRead(t *testing.T) {
	Convey("Given a gmail sender and missing report file", t, func() {
		sender := &gmailWorkoutReportSender{
			httpClient:    &http.Client{},
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        "http://localhost",
		}

		Convey("When Send is called", func() {
			err := sender.Send(filepath.Join(t.TempDir(), "missing.xlsx"))

			Convey("Then it should return read report file error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to read report file")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorWhenRequestCreationFails(t *testing.T) {
	Convey("Given a gmail sender with invalid api URL", t, func() {
		sender := &gmailWorkoutReportSender{
			httpClient:    &http.Client{},
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        "://invalid-url",
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should return request creation error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to create gmail request")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorWhenHTTPClientDoFails(t *testing.T) {
	Convey("Given a gmail sender whose http client fails", t, func() {
		httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("dial failure")
		})}
		sender := &gmailWorkoutReportSender{
			httpClient:    httpClient,
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        "https://gmail.googleapis.com/gmail/v1/users/me/messages/send",
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should return send request error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to send gmail request")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorWhenResponseDecodeFails(t *testing.T) {
	Convey("Given a gmail sender and malformed json response", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"id":`))
		}))
		defer server.Close()

		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        server.URL,
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it should return decode response error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to decode gmail response")
			})
		})
	})
}

func TestWrapBase64ShouldReturnEmptyStringWhenInputIsEmpty(t *testing.T) {
	Convey("Given an empty base64 string", t, func() {
		Convey("When wrapBase64 is called", func() {
			result := wrapBase64("")

			Convey("Then it should return empty string", func() {
				So(result, ShouldEqual, "")
			})
		})
	})
}

func TestWrapBase64ShouldInsertCRLFWhenInputExceedsLineLength(t *testing.T) {
	Convey("Given a base64 string longer than one line", t, func() {
		value := strings.Repeat("A", 80)

		Convey("When wrapBase64 is called", func() {
			result := wrapBase64(value)

			Convey("Then it should break lines using CRLF", func() {
				So(result, ShouldContainSubstring, "\r\n")
				So(strings.HasSuffix(result, "\r\n"), ShouldBeTrue)
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorOnUnauthorizedResponse(t *testing.T) {
	Convey("Given the Gmail API returns unauthorized", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "denied", http.StatusUnauthorized)
		}))
		defer server.Close()

		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        server.URL,
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it returns the Gmail API status error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "gmail api error: 401 Unauthorized")
			})
		})
	})
}

func TestGmailWorkoutReportSenderSendReturnsErrorOnMalformedResponse(t *testing.T) {
	Convey("Given the Gmail API returns a malformed success response", t, func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"threadId":"missing-id"}`))
		}))
		defer server.Close()

		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
			from:          "from@example.com",
			to:            "to@example.com",
			subject:       "Cycling Workout Report",
			apiURL:        server.URL,
		}

		Convey("When Send is called", func() {
			err := sender.Send(createTempReportFile(t))

			Convey("Then it reports the malformed Gmail response", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "malformed gmail response")
			})
		})
	})
}

func createTempReportFile(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "report.xlsx")
	err := os.WriteFile(filePath, []byte(strings.Repeat("xlsx-bytes", 10)), 0o600)
	if err != nil {
		t.Fatalf("failed to create temp report file: %v", err)
	}
	return filePath
}
