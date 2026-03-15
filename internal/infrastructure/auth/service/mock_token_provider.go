package auth_service

import (
	authInterfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type mockTokenProvider struct {
	Tokens *auth_model.Tokens
	Err    error
}

// RefreshGoogleToken implements [auth_interfaces.OAuthProvider].
func (m *mockTokenProvider) RefreshGoogleToken(refreshToken string) (*auth_model.Token, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Tokens == nil {
		return nil, nil
	}
	return m.Tokens.GoogleToken, m.Err
}

// RefreshStravaToken implements [authInterfaces.OAuthProvider].
func (m *mockTokenProvider) RefreshStravaToken(refreshToken string) (*auth_model.Token, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Tokens == nil {
		return nil, nil
	}
	return m.Tokens.StravaToken, m.Err
}

var _ authInterfaces.OAuthProvider = (*mockTokenProvider)(nil)
