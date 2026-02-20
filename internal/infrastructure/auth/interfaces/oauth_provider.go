package auth_interfaces

import auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type OAuthProvider interface {
	RefreshToken(refreshToken string) (*auth_model.Token, error)
}
