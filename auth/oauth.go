package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// GetAuthorizationCode starts a local server and captures the OAuth code
func (c *Client) GetAuthorizationCode() (string, error) {
	codeChan := make(chan string)
	errChan := make(chan error)

	// Create HTTP server to catch callback
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No authorization code received", http.StatusBadRequest)
			errChan <- fmt.Errorf("no code in callback")
			return
		}

		// Send success response to browser
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `
			<html>
			<body style="font-family: Arial; text-align: center; padding: 50px;">
				<h1>Authorization Successful</h1>
				<p>You can close this window and return to your application.</p>
			</body>
			</html>
		`) //nolint:errcheck // Error writing to response writer is non-critical

		codeChan <- code
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Give server time to start
	time.Sleep(500 * time.Millisecond)

	// Build authorization URL
	authURL := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=activity:read_all",
		c.Config.ClientID, c.Config.RedirectURI,
	)

	fmt.Println("\nPlease open this URL in your browser to authorize:")
	fmt.Println(authURL)
	fmt.Println("\nWaiting for authorization callback...")

	// Wait for code or timeout
	select {
	case code := <-codeChan:
		// Shutdown server gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
		return code, nil

	case err := <-errChan:
		_ = server.Shutdown(context.Background())
		return "", err

	case <-time.After(5 * time.Minute):
		_ = server.Shutdown(context.Background())
		return "", fmt.Errorf("timeout waiting for authorization (5 minutes)")
	}
}
