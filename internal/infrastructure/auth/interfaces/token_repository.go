package auth_interfaces

import auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type TokenRepository interface {
	GetTokens() (*auth_model.Token, error)
	SaveTokens(token *auth_model.Token) error
}
