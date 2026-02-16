package interfaces

import "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"

type TokenRepository interface {
	GetTokens() (*model.Token, error)
	SaveTokens(token *model.Token) error
}
