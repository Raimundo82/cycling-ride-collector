package token_provider

import (
	"context"

	token_client "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/client"
	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type TokenProvider interface {
	GetValidToken(ctx context.Context) (string, error)
}

type tokenProvider struct {
	RefreshTokenInput *token_model.RefreshTokenInput
	TokenClient       token_client.TokenClient
}

func NewTokenProvider(input *token_model.RefreshTokenInput, tokenClient token_client.TokenClient) TokenProvider {
	return &tokenProvider{
		RefreshTokenInput: input,
		TokenClient:       tokenClient,
	}
}

// GetToken implements [TokenProvider].
func (t *tokenProvider) GetValidToken(ctx context.Context) (string, error) {
	resp, err := t.TokenClient.RefreshToken(ctx, t.RefreshTokenInput)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

var _ TokenProvider = (*tokenProvider)(nil)
