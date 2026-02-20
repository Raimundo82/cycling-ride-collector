package auth_repository

import (
	"encoding/json"
	"fmt"
	"os"

	auth_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type tokenRepository struct {
	token    *auth_model.Token
	filePath string
}

var _ auth_interfaces.TokenRepository = (*tokenRepository)(nil)

func NewTokenRepository(filePath string) (*tokenRepository, error) {
	repo := &tokenRepository{filePath: filePath}
	tokens, err := repo.GetTokens()
	if err != nil {
		return nil, err
	}
	repo.token = tokens

	return repo, nil
}

// GetTokens implements [auth_interfaces.TokenRepository].
func (t *tokenRepository) GetTokens() (*auth_model.Token, error) {
	if t.token != nil {
		return t.token, nil
	}

	file, err := os.Open(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	defer func() { _ = file.Close() }()

	var token auth_model.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}
	return &token, nil
}

// SaveTokens implements [auth_interfaces.TokenRepository].
func (t *tokenRepository) SaveTokens(token *auth_model.Token) error {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to serialize token: %w", err)
	}

	err = os.WriteFile(t.filePath, tokenBytes, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	err = os.Chmod(t.filePath, 0o600)
	if err != nil {
		return fmt.Errorf("failed to set token file permissions: %w", err)
	}
	t.token = token
	return nil
}
