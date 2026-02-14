package usecase

import (
	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
)

type GetAccessToken struct {
	TokenProvider   contracts.TokenProvider
	TokenRepository contracts.TokenRepository
}

func (uc *GetAccessToken) Execute() (*dto.Token, error) {
	token, err := uc.TokenRepository.GetTokens()
	if err != nil {
		return nil, err
	}

	if token == nil || token.IsExpired() {
		token, err = uc.TokenProvider.RefreshAccessToken()
		if err != nil {
			return nil, err
		}
		err = uc.TokenRepository.SaveTokens(token)
		if err != nil {
			return nil, err
		}
	}
	return token, nil
}
