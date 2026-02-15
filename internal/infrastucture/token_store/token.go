package tokens

import "time"

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-1 * time.Minute))
}
