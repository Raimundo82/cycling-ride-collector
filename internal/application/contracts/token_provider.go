package contracts

import "github.com/raimundo82/cycling-ride-collector/internal/application/dto"

type TokenProvider interface {
	RefreshAccessToken() (*dto.Token, error)
}
