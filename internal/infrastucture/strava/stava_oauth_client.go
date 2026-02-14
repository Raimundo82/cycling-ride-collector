package strava

type StravaOAuthClient interface {
	RefreshAccessToken(req *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
}
