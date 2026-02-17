package provider

import "context"

type TokenRefresher interface {
	RefreshToken(context.Context, *RefreshAccessTokenRequest) (*RefreshAccessTokenResponse, error)
}
