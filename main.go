package main

import (
	"fmt"
	"log"
	"time"

	"github.com/raimundo82/go-strava-weekly/auth"
)

func main() {
	fmt.Println("== Strava OAuth Token Fetcher ==")
	config := auth.LoadConfig()

	// Create Strava client with dependency injection
	client := auth.NewClient(config)

	// Get authorization code via OAuth flow
	code, err := client.GetAuthorizationCode()
	if err != nil {
		log.Fatalf("Failed to get authorization: %v", err)
	}

	fmt.Println("\n✓ Authorization code received")
	fmt.Println("Exchanging code for access token...")

	// Exchange code for access token
	token, err := client.GetToken(code)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	fmt.Println("\n=== Access Token ===")
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Printf("Refresh Token: %s\n", token.RefreshToken)
	fmt.Printf("Athlete ID: %d\n", token.Athlete.ID)
	fmt.Printf("Expires At: %s\n", time.Unix(token.ExpiresAt, 0).Format(time.RFC3339))
	fmt.Printf("Athlete: %s %s (@%s)\n", token.Athlete.FirstName, token.Athlete.LastName, token.Athlete.Username)
	fmt.Println("====================")
}
