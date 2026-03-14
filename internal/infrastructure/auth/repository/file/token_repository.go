package auth_repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	auth_interfaces "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/interfaces"
	auth_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type tokenRepository struct {
	tokens   *auth_model.Tokens
	filePath string
}

var _ auth_interfaces.TokenRepository = (*tokenRepository)(nil)

func NewTokenRepository(filePath string) (*tokenRepository, error) {
	repo := &tokenRepository{filePath: filePath}
	tokens, err := repo.GetTokens()
	if err != nil {
		return nil, err
	}

	repo.tokens = ensureTokens(tokens)
	return repo, nil
}

// GetTokens implements [auth_interfaces.TokenRepository].
func (t *tokenRepository) GetTokens() (*auth_model.Tokens, error) {
	if t.tokens != nil {
		return t.tokens, nil
	}

	tokens, err := t.readTokensFromFile()
	if err != nil {
		return nil, err
	}
	t.tokens = tokens
	return t.tokens, nil
}

func (t *tokenRepository) readTokensFromFile() (*auth_model.Tokens, error) {
	fileBytes, err := os.ReadFile(t.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}

	var tokens auth_model.Tokens
	if err := json.Unmarshal(fileBytes, &tokens); err != nil {
		return nil, fmt.Errorf("failed to decode tokens: %w", err)
	}
	return &tokens, nil
}

func ensureTokens(tokens *auth_model.Tokens) *auth_model.Tokens {
	if tokens == nil {
		tokens = &auth_model.Tokens{}
	}

	if tokens.StravaToken == nil {
		tokens.StravaToken = emptyToken()
	}

	if tokens.GoogleToken == nil {
		tokens.GoogleToken = emptyToken()
	}

	return tokens
}

func emptyToken() *auth_model.Token {
	return &auth_model.Token{
		AccessToken:  "",
		RefreshToken: "",
		ExpiresAt:    time.Time{},
	}
}

// SaveStravaToken implements [auth_interfaces.TokenRepository].
func (t *tokenRepository) SaveStravaToken(token *auth_model.Token) error {
	currentTokens := ensureTokens(t.tokens)
	tokens := &auth_model.Tokens{
		StravaToken: token,
		GoogleToken: currentTokens.GoogleToken,
	}

	tokenBytes, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to serialize tokens: %w", err)
	}

	err = writeTokenFileAtomically(t.filePath, tokenBytes)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	t.tokens = tokens
	return nil
}

// SaveGoogleToken implements [auth_interfaces.TokenRepository].
func (t *tokenRepository) SaveGoogleToken(token *auth_model.Token) error {
	currentTokens := ensureTokens(t.tokens)
	tokens := &auth_model.Tokens{
		GoogleToken: token,
		StravaToken: currentTokens.StravaToken,
	}

	tokenBytes, err := json.Marshal(tokens)
	if err != nil {
		return fmt.Errorf("failed to serialize tokens: %w", err)
	}

	err = writeTokenFileAtomically(t.filePath, tokenBytes)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	t.tokens = tokens
	return nil
}

func writeTokenFileAtomically(filePath string, content []byte) error {
	dir := filepath.Dir(filePath)
	if dir == "" {
		dir = "."
	}

	tempFile, err := os.CreateTemp(dir, ".token-*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() { _ = os.Remove(tempPath) }()

	if err := tempFile.Chmod(0o600); err != nil {
		_ = tempFile.Close()
		return err
	}
	if _, err := tempFile.Write(content); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, filePath); err != nil {
		return err
	}
	return os.Chmod(filePath, 0o600)
}
