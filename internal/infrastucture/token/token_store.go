package token

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/raimundo82/cycling-ride-collector/internal/application/contracts"
	"github.com/raimundo82/cycling-ride-collector/internal/application/dto"
)

type tokenStore struct {
	token    dto.Token
	filePath string
}

// GetTokens implements [contracts.TokenRepository].
func (t *tokenStore) GetTokens() (*dto.Token, error) {
	file, err := os.Open(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	defer func() { _ = file.Close() }()

	var store tokenStore
	if err := json.NewDecoder(file).Decode(&store.token); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}
	return &store.token, nil
}

// SaveTokens implements [contracts.TokenRepository].
func (t *tokenStore) SaveTokens(token *dto.Token) error {
	file, err := os.Create(t.filePath)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer func() { _ = file.Close() }()

	return json.NewEncoder(file).Encode(token)
}

var _ contracts.TokenRepository = (*tokenStore)(nil)
