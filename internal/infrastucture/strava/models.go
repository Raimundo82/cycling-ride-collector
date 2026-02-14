package strava

type ActivityDto struct {
	ID                 int64   `json:"id"`
	SportType          string  `json:"sport_type"`
	Commute            bool    `json:"commute"`
	IsTrainer          bool    `json:"trainer"`
	WorkoutType        int     `json:"workout_type"`
	StartDate          string  `json:"start_date_local"`
	Distance           float64 `json:"distance"`
	Duration           int     `json:"moving_time"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
	AveragePower       float64 `json:"average_watts"`
	HasHeartRate       bool    `json:"has_heartrate"`
	DeviceWatts        bool    `json:"device_watts"`
	WeightedAvgPower   float64 `json:"weighted_average_watts"`
	AverageHeartRate   float64 `json:"average_heartrate"`
	MaxHeartRate       float64 `json:"max_heartrate"`
	AverageCadence     float64 `json:"average_cadence"`
	LegSensations      string
	Watts              *WattsStreamDto
}

type DetailedActivityDto struct {
	ID            int64  `json:"id"`
	LegSensations string `json:"private_note"`
}

type WattsStreamDto struct {
	WattsData []int `json:"data"`
}

type wattsStreamResponse struct {
	Watts WattsStreamDto `json:"watts"`
}

type RefreshAccessTokenResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}
