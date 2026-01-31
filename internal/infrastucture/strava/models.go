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
	WeightedAvgPower   float64 `json:"weighted_average_watts"`
	AverageHeartRate   float64 `json:"average_heartrate"`
	MaxHeartRate       float64 `json:"max_heartrate"`
	AverageCadence     float64 `json:"average_cadence"`
	Watts              *WattsStreamDto
}

type WattsStreamDto struct {
	WattsData []int `json:"data"`
}

type wattsStreamResponse struct {
	Watts WattsStreamDto `json:"watts"`
}
