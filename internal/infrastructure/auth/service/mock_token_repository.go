package auth_service

import (
	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type mockTokenRepository struct {
	SaveErr error
	GetErr  error
	Tokens  *auth_model.Tokens
}

// GetTokens implements [authInterfaces.TokenRepository].
func (m *mockTokenRepository) GetTokens() (*auth_model.Tokens, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Tokens == nil {
		return &auth_model.Tokens{}, nil
	}
	return m.Tokens, m.GetErr
}

// SaveStravaToken implements [authInterfaces.TokenRepository].
func (m *mockTokenRepository) SaveStravaToken(token *auth_model.Token) error {
	return m.SaveErr
}

// SaveGoogleToken implements [authInterfaces.TokenRepository].
func (m *mockTokenRepository) SaveGoogleToken(token *auth_model.Token) error {
	return m.SaveErr
}

var _ authInterfaces.TokenRepository = (*mockTokenRepository)(nil)
