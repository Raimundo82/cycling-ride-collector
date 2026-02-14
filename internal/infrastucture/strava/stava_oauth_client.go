package strava

// NOTE: This file is currently named `stava_oauth_client.go`; it should be renamed
// to `strava_oauth_client.go` for consistency with the package and type names.
type StravaOAuthClient interface {
	RefreshAccessToken(req *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
}
