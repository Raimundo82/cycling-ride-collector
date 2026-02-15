package strava

type AccessTokenRefresher interface {
	RefreshAccessToken(req *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
}
