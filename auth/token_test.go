package auth

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/raimundo82/go-strava-weekly/models"
)

// mockHTTPClient is a mock implementation of HTTPClient for testing
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return m.response, m.err
}

func TestGetToken_Success(t *testing.T) {
	// Mock successful response
	mockResponse := `{
		"token_type": "Bearer",
		"expires_at": 1700000000,
		"expires_in": 21600,
		"refresh_token": "test_refresh_token",
		"access_token": "test_access_token",
		"athlete": {
			"id": 12345,
			"username": "testuser",
			"firstname": "John",
			"lastname": "Doe"
		}
	}`

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(mockResponse)),
		},
		err: nil,
	}

	config := Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		BaseURL:      "https://test.strava.com",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client := NewClientWithHTTP(config, mockClient)

	token, err := client.GetToken("test_code")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token.AccessToken != "test_access_token" {
		t.Errorf("Expected access token 'test_access_token', got '%s'", token.AccessToken)
	}

	if token.RefreshToken != "test_refresh_token" {
		t.Errorf("Expected refresh token 'test_refresh_token', got '%s'", token.RefreshToken)
	}

	if token.Athlete.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", token.Athlete.Username)
	}
}

func TestGetToken_BadRequest(t *testing.T) {
	mockErrorResponse := `{"message":"Bad Request","errors":[{"resource":"RequestToken","field":"code","code":"invalid"}]}`

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(mockErrorResponse)),
		},
		err: nil,
	}

	config := Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		BaseURL:      "https://test.strava.com",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client := NewClientWithHTTP(config, mockClient)

	_, err := client.GetToken("invalid_code")
	if err == nil {
		t.Fatal("Expected error for bad request, got nil")
	}

	if !strings.Contains(err.Error(), "token request failed") {
		t.Errorf("Expected error message to contain 'token request failed', got: %v", err)
	}
}

func TestGetToken_InvalidJSON(t *testing.T) {
	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("invalid json")),
		},
		err: nil,
	}

	config := Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		BaseURL:      "https://test.strava.com",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client := NewClientWithHTTP(config, mockClient)

	_, err := client.GetToken("test_code")
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse token response") {
		t.Errorf("Expected error message to contain 'failed to parse token response', got: %v", err)
	}
}

func TestGetToken_NetworkError(t *testing.T) {
	mockClient := &mockHTTPClient{
		response: nil,
		err:      io.EOF,
	}

	config := Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		BaseURL:      "https://test.strava.com",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client := NewClientWithHTTP(config, mockClient)

	_, err := client.GetToken("test_code")
	if err == nil {
		t.Fatal("Expected network error, got nil")
	}

	if !strings.Contains(err.Error(), "failed to request token") {
		t.Errorf("Expected error message to contain 'failed to request token', got: %v", err)
	}
}

func TestNewClient(t *testing.T) {
	config := Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		BaseURL:      "https://test.strava.com",
		RedirectURI:  "http://localhost:8080/callback",
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	if client.Config.ClientID != "test_id" {
		t.Errorf("Expected client ID 'test_id', got '%s'", client.Config.ClientID)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized, got nil")
	}
}

// Example of table-driven tests
func TestGetToken_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedError bool
		errorContains string
		expectedToken *models.TokenResponse
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			responseBody: `{
				"token_type": "Bearer",
				"expires_at": 1700000000,
				"expires_in": 21600,
				"refresh_token": "refresh123",
				"access_token": "access123",
				"athlete": {"id": 1, "username": "test"}
			}`,
			expectedError: false,
			expectedToken: &models.TokenResponse{
				TokenType:    "Bearer",
				AccessToken:  "access123",
				RefreshToken: "refresh123",
			},
		},
		{
			name:          "Unauthorized",
			statusCode:    http.StatusUnauthorized,
			responseBody:  `{"message":"Unauthorized"}`,
			expectedError: true,
			errorContains: "token request failed",
		},
		{
			name:          "Invalid JSON",
			statusCode:    http.StatusOK,
			responseBody:  `not json`,
			expectedError: true,
			errorContains: "failed to parse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				response: &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
				},
			}

			config := Config{
				ClientID:     "test",
				ClientSecret: "secret",
				BaseURL:      "https://test.com",
			}

			client := NewClientWithHTTP(config, mockClient)
			token, err := client.GetToken("code")

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Error should contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if token.AccessToken != tt.expectedToken.AccessToken {
					t.Errorf("Expected token %s, got %s", tt.expectedToken.AccessToken, token.AccessToken)
				}
			}
		})
	}
}
