package auth_service

import (
	"fmt"
	"time"

	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type GoogleTokenService struct {
	oauthProvider authInterfaces.OAuthProvider
	tokenRepo     authInterfaces.TokenRepository
	cachedToken   *auth_model.Token
}

func NewGoogleTokenService(
	oauthProvider authInterfaces.OAuthProvider,
	tokenRepo authInterfaces.TokenRepository,
) *GoogleTokenService {
	return &GoogleTokenService{
		oauthProvider: oauthProvider,
		tokenRepo:     tokenRepo,
	}
}

var _ authInterfaces.TokenProvider = (*GoogleTokenService)(nil)

func (s *GoogleTokenService) GetValidToken() (string, error) {
	if s.cachedToken != nil && s.cachedToken.AccessToken != "" && !s.cachedToken.IsExpired() {
		return s.cachedToken.AccessToken, nil
	}

	tokens, err := s.tokenRepo.GetTokens()
	if err != nil {
		return "", err
	}
	token := tokens.GoogleToken

	if token == nil {
		return "", fmt.Errorf("no google tokens available")
	}

	if token.RefreshToken == "" {
		return "", fmt.Errorf("no google refresh token available")
	}

	token, err = s.oauthProvider.RefreshGoogleToken(token.RefreshToken)
	if err != nil {
		return "", err
	}

	token.RefreshToken = tokens.GoogleToken.RefreshToken
	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	s.cachedToken = token

	if err := s.tokenRepo.SaveGoogleToken(token); err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
