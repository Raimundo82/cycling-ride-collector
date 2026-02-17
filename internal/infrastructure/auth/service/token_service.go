package auth

import (
	"fmt"

	activityInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/activity/interfaces"
	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
)

type tokenService struct {
	oauthProvider authInterfaces.OAuthProvider
	tokenRepo     authInterfaces.TokenRepository
}

func NewTokenService(oauthProvider authInterfaces.OAuthProvider, tokenRepo authInterfaces.TokenRepository) *tokenService {
	return &tokenService{oauthProvider: oauthProvider, tokenRepo: tokenRepo}
}

var _ activityInterfaces.TokenProvider = (*tokenService)(nil)

func (s *tokenService) GetValidToken() (string, error) {
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
		tokens, err = s.oauthProvider.RefreshToken(tokens.RefreshToken)
		if err != nil {
			return "", err
		}
		if err := s.tokenRepo.SaveTokens(tokens); err != nil {
			return "", err
		}
	}
	return tokens.AccessToken, nil
}
