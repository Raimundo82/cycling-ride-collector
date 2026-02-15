package strava

import (
	"fmt"

	tokens "github.com/raimundo82/cycling-ride-collector/internal/infrastucture/token_store"
)

type tokenService struct {
	oauthProvider OAuthProvider
	tokenRepo     tokens.TokenRepository
}

func NewTokenService(oauthProvider OAuthProvider, tokenRepo tokens.TokenRepository) *tokenService {
	return &tokenService{oauthProvider: oauthProvider, tokenRepo: tokenRepo}
}

func (s *tokenService) GetValidAccessToken() (string, error) {
	tokens, err := s.tokenRepo.GetTokens()
	if err != nil {
		return "", err
	}

	if tokens == nil {
		return "", fmt.Errorf("no tokens available")
	}

	if tokens.IsExpired() || tokens.AccessToken == "" {
		if tokens.RefreshToken == "" {
			return "", fmt.Errorf("no valid access token and no refresh token available")
		}
		tokens, err = s.oauthProvider.RefreshAccessToken(tokens.RefreshToken)
		if err != nil {
			return "", err
		}
		if err := s.tokenRepo.SaveTokens(tokens); err != nil {
			return "", err
		}
	}
	return tokens.AccessToken, nil
}
