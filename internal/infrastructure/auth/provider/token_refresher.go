package auth_provider

import (
	"context"

	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type TokenRefresher interface {
	RefreshToken(context.Context, *auth_model.RefreshAccessTokenRequest) (*auth_model.RefreshAccessTokenResponse, error)
}
