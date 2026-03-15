package auth_service

import (
	"fmt"

	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type StravaTokenService struct {
	oauthProvider authInterfaces.OAuthProvider
	tokenRepo     authInterfaces.TokenRepository
	cachedToken   *auth_model.Token
}

func NewStravaTokenService(
	oauthProvider authInterfaces.OAuthProvider,
	tokenRepo authInterfaces.TokenRepository,
) *StravaTokenService {
	return &StravaTokenService{
		oauthProvider: oauthProvider,
		tokenRepo:     tokenRepo,
	}
}

var _ authInterfaces.TokenProvider = (*StravaTokenService)(nil)

func (s *StravaTokenService) GetValidToken() (string, error) {
	if s.cachedToken != nil && s.cachedToken.AccessToken != "" && !s.cachedToken.IsExpired() {
		return s.cachedToken.AccessToken, nil
	}

	tokens, err := s.tokenRepo.GetTokens()
	if err != nil {
		return "", err
	}

	token := tokens.StravaToken
	if token == nil {
		return "", fmt.Errorf("no strava tokens available")
	}

	if token.IsExpired() || token.AccessToken == "" {
		if token.RefreshToken == "" {
			return "", fmt.Errorf("no strava refresh token available")
		}

		token, err = s.oauthProvider.RefreshStravaToken(token.RefreshToken)
		if err != nil {
			return "", err
		}
		s.cachedToken = token
		if err := s.tokenRepo.SaveStravaToken(token); err != nil {
			return "", err
		}
	}
	return token.AccessToken, nil
}
