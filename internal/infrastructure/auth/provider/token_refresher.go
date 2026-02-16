package provider

type TokenRefresher interface {
	RefreshToken(*RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
}
