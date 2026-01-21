package auth

import (
	"io"
	"net/http"
)

// HTTPClient interface allows for mocking HTTP calls in tests
type HTTPClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

// Client represents a Strava API client with configurable dependencies
type Client struct {
	HTTPClient HTTPClient
	Config     Config
}

// NewClient creates a new Strava API client with the given configuration
func NewClient(config Config) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		Config:     config,
	}
}

// NewClientWithHTTP creates a new client with a custom HTTP client (useful for testing)
func NewClientWithHTTP(config Config, httpClient HTTPClient) *Client {
	return &Client{
		HTTPClient: httpClient,
		Config:     config,
	}
}
