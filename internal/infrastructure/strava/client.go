package strava

import (
	"fmt"
	"time"

	"github.com/raimundo82/go-strava-weekly/internal/domain"
)

// Client implements the ActivitiesRepository interface for Strava API
type Client struct {
	apiKey string
	// In a real implementation, this would contain HTTP client, auth tokens, etc.
}

// NewClient creates a new Strava API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// GetActivitiesBetween fetches activities from Strava between the given dates
func (c *Client) GetActivitiesBetween(startDate, endDate time.Time) ([]domain.Activity, error) {
	// This is a placeholder implementation
	// In a real implementation, this would:
	// 1. Make HTTP requests to Strava API
	// 2. Handle authentication and rate limiting
	// 3. Parse the response
	// 4. Map Strava API response to domain.Activity

	fmt.Println("Step 1: Calling Strava API (placeholder implementation)")
	fmt.Printf("  → Would fetch activities from %s to %s\n",
		startDate.Format(time.RFC3339),
		endDate.Format(time.RFC3339))

	// Return empty slice for now - ready for real implementation
	return []domain.Activity{}, nil
}
