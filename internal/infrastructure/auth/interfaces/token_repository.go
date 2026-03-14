package auth_interfaces

import auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type TokenRepository interface {
	GetTokens() (*auth_model.Tokens, error)
	SaveStravaToken(token *auth_model.Token) error
	SaveGoogleToken(token *auth_model.Token) error
}
