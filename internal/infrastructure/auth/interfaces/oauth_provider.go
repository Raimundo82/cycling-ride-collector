package auth_interfaces

import auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type OAuthProvider interface {
	RefreshStravaToken(refreshToken string) (*auth_model.Token, error)
	RefreshGoogleToken(refreshToken string) (*auth_model.Token, error)
}
