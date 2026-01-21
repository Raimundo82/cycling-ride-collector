package models

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
	} `json:"athlete"`
}

type Activity struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Distance           float64 `json:"distance"`
	MovingTime         int     `json:"moving_time"`
	ElapsedTime        int     `json:"elapsed_time"`
	Type               string  `json:"type"`
	StartDate          string  `json:"start_date"`
	StartDateLocal     string  `json:"start_date_local"`
	AverageSpeed       float64 `json:"average_speed"`
	MaxSpeed           float64 `json:"max_speed"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
}
