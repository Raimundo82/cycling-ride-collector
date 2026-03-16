package token_model

import "time"

type Token struct {
	accessToken string
	expiresAt   time.Time
}

func NewToken(accessToken string, expiresIn int) Token {
	return Token{
		accessToken: accessToken,
		expiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}

func (t *Token) IsValid() bool {
	return t.accessToken != "" && time.Now().Add(30*time.Second).Before(t.expiresAt)
}

func (t *Token) AccessToken() string {
	return t.accessToken
}

func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}

type RefreshTokenInput struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutput struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
