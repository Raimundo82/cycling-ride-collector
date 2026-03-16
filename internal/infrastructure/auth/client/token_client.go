package token_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	token_model "github.com/raimundo82/cycling-ride-collector/internal/infrastructure/auth/model"
)

type TokenClient interface {
	RefreshToken(ctx context.Context, input *token_model.RefreshTokenInput) (*token_model.RefreshTokenOutput, error)
}

type tokenClient struct {
	tokenUrl   string
	httpClient *http.Client
}

func NewTokenClient(tokenUrl string) TokenClient {
	return &tokenClient{
		tokenUrl:   tokenUrl,
		httpClient: http.DefaultClient,
	}
}
// RefreshToken implements [TokenProvider].
func (o *tokenClient) RefreshToken(ctx context.Context, input *token_model.RefreshTokenInput) (*token_model.RefreshTokenOutput, error) {
	form := url.Values{}
	form.Set("client_id", input.ClientID)
	form.Set("client_secret", input.ClientSecret)
	form.Set("grant_type", input.GrantType)
	form.Set("refresh_token", input.RefreshToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google error: %s", resp.Status)
	}

	var refreshResponse token_model.RefreshTokenOutput
	if err := json.NewDecoder(resp.Body).Decode(&refreshResponse); err != nil {
		return nil, err
	}
	return &refreshResponse, nil
}
var _ TokenClient = (*tokenClient)(nil)
