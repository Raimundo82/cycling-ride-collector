package email

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/config"
	auth_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
)

const (
	gmailSendMessageURL = "https://gmail.googleapis.com/gmail/v1/users/me/messages/send"
	contentTypeHeader   = "Content-Type"
)

type gmailWorkoutReportSender struct {
	httpClient    *http.Client
	tokenProvider auth_interfaces.TokenProvider
	from          string
	to            string
	subject       string
	apiURL        string
}

type gmailSendRequest struct {
	Raw string `json:"raw"`
}

type gmailSendResponse struct {
	ID string `json:"id"`
}

var _ contracts.WorkoutReportSender = (*gmailWorkoutReportSender)(nil)

func NewGmailWorkoutReportSender(httpClient *http.Client, tokenProvider auth_interfaces.TokenProvider, cfg *config.Config) *gmailWorkoutReportSender {
	return &gmailWorkoutReportSender{
		httpClient:    httpClient,
		tokenProvider: tokenProvider,
		from:          cfg.Email.From,
		to:            cfg.Email.To,
		subject:       cfg.Email.Subject,
		apiURL:        gmailSendMessageURL,
	}
}

func (g *gmailWorkoutReportSender) Send(reportPath string) error {
	accessToken, err := g.tokenProvider.GetValidToken()
	if err != nil {
		return fmt.Errorf("failed to get google access token: %w", err)
	}

	rawMessage, err := g.buildMessage(reportPath)
	if err != nil {
		return err
	}

	body, err := json.Marshal(gmailSendRequest{Raw: base64.RawURLEncoding.EncodeToString(rawMessage)})
	if err != nil {
		return fmt.Errorf("failed to encode gmail request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, g.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create gmail request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set(contentTypeHeader, "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send gmail request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("gmail api error: %s", resp.Status)
	}

	var sendResponse gmailSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&sendResponse); err != nil {
		return fmt.Errorf("failed to decode gmail response: %w", err)
	}
	if sendResponse.ID == "" {
		return fmt.Errorf("malformed gmail response: missing message id")
	}

	return nil
}

func (g *gmailWorkoutReportSender) buildMessage(reportPath string) ([]byte, error) {
	attachment, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	textHeader := textproto.MIMEHeader{}
	textHeader.Set(contentTypeHeader, "text/plain; charset=UTF-8")
	textPart, err := writer.CreatePart(textHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create email body part: %w", err)
	}
	if _, err := textPart.Write([]byte("Cycling workout report attached.")); err != nil {
		return nil, fmt.Errorf("failed to write email body: %w", err)
	}

	attachmentHeader := textproto.MIMEHeader{}
	attachmentHeader.Set(contentTypeHeader, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	attachmentHeader.Set("Content-Transfer-Encoding", "base64")
	attachmentHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(reportPath)))
	attachmentPart, err := writer.CreatePart(attachmentHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment part: %w", err)
	}
	if _, err := attachmentPart.Write([]byte(wrapBase64(base64.StdEncoding.EncodeToString(attachment)))); err != nil {
		return nil, fmt.Errorf("failed to write attachment body: %w", err)
	}

	boundary := writer.Boundary()
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize email body: %w", err)
	}

	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %s\r\n", g.from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", g.to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", g.subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	message.WriteString("\r\n")
	message.Write(body.Bytes())

	return message.Bytes(), nil
}

func wrapBase64(value string) string {
	if value == "" {
		return ""
	}

	const lineLength = 76
	var builder strings.Builder
	for start := 0; start < len(value); start += lineLength {
		end := start + lineLength
		if end > len(value) {
			end = len(value)
		}
		builder.WriteString(value[start:end])
		builder.WriteString("\r\n")
	}
	return builder.String()
}
