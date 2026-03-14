package email

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	auth_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	. "github.com/smartystreets/goconvey/convey"
)

type stubGoogleTokenProvider struct {
	Token string
	Err   error
}

var _ auth_interfaces.TokenProvider = (*stubGoogleTokenProvider)(nil)

func (s *stubGoogleTokenProvider) GetValidToken() (string, error) {
	return s.Token, s.Err
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
		sender := &gmailWorkoutReportSender{
			httpClient:    server.Client(),
			tokenProvider: &stubGoogleTokenProvider{Token: "google-access-token"},
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

				decodedMessage, decodeErr := base64.RawURLEncoding.DecodeString(rawMessage)
				So(decodeErr, ShouldBeNil)
				So(string(decodedMessage), ShouldContainSubstring, "Subject: Cycling Workout Report")
				So(string(decodedMessage), ShouldContainSubstring, "filename=\"report.xlsx\"")
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
