package tokens

import (
	"encoding/json"
	"fmt"
	"os"
)

type TokenRepository interface {
	GetTokens() (*Token, error)
	SaveTokens(token *Token) error
}

type TokenStore struct {
	token    Token
	filePath string
}

func NewTokenStore(filePath string) TokenRepository {
	return &TokenStore{filePath: filePath}
}

// GetTokens implements [contracts.TokenRepository].
func (t *TokenStore) GetTokens() (*Token, error) {
	file, err := os.Open(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	defer func() { _ = file.Close() }()

	var store TokenStore
	if err := json.NewDecoder(file).Decode(&store.token); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}
	return &store.token, nil
}

// SaveTokens implements [contracts.TokenRepository].
func (t *TokenStore) SaveTokens(token *Token) error {
	file, err := os.Create(t.filePath)
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer func() { _ = file.Close() }()

	return json.NewEncoder(file).Encode(token)
}
