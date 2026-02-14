package contracts

import (
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
)

type TokenRepository interface {
	SaveTokens(token *dto.Token) error
	GetTokens() (*dto.Token, error)
}
