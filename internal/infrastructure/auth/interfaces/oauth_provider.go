package interfaces

import "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type OAuthProvider interface {
	RefreshToken(refreshToken string) (*model.Token, error)
}
