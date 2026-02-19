package model

type DetailedAthlete struct {
	ID     int64   `json:"id"`
	Weight float64 `json:"weight"`
	Zones  Zones
}
