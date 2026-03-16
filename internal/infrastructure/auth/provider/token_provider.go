package token_provider

import (
	"context"
	"time"

	token_client "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/client"
	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type TokenProvider interface {
	GetValidToken(ctx context.Context) (string, error)
}

const expiryBuffer = 30 * time.Second

type tokenProvider struct {
	RefreshTokenInput *token_model.RefreshTokenInput
	TokenClient       token_client.TokenClient
	cachedToken       string
	tokenExpiresAt    time.Time
}

func NewTokenProvider(input *token_model.RefreshTokenInput, tokenClient token_client.TokenClient) TokenProvider {
	return &tokenProvider{
		RefreshTokenInput: input,
		TokenClient:       tokenClient,
	}
}

// GetValidToken implements [TokenProvider].
func (t *tokenProvider) GetValidToken(ctx context.Context) (string, error) {
	if t.cachedToken != "" && time.Now().Add(expiryBuffer).Before(t.tokenExpiresAt) {
		return t.cachedToken, nil
	}

	resp, err := t.TokenClient.RefreshToken(ctx, t.RefreshTokenInput)
	if err != nil {
		return "", err
	}

	expiresIn := resp.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	t.cachedToken = resp.AccessToken
	t.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return t.cachedToken, nil
}

var _ TokenProvider = (*tokenProvider)(nil)
