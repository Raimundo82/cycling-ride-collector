package auth_model

import "time"

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int64     `json:"expires_in,omitempty"`
}

type Tokens struct {
	StravaToken *Token `json:"strava"`
	GoogleToken *Token `json:"google"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-1 * time.Minute))
}
